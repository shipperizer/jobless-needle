package tasker

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type ListPermissionsResultValue struct {
	Permissions []string
	// token       string
	Err error
}

func NewListPermissionsResultValue(permissions []string, err error) *ListPermissionsResultValue {
	v := new(ListPermissionsResultValue)
	v.Permissions = permissions
	v.Err = err

	return v
}

type ListPermissionsResult struct {
	key   uuid.UUID
	value ListPermissionsResultValue
}

func (r *ListPermissionsResult) GetKey() string {
	return r.key.String()
}

func (r *ListPermissionsResult) GetValue() interface{} {
	return r.value
}

func (r *ListPermissionsResult) SetValue(v interface{}) {

	value, ok := v.(ListPermissionsResultValue)

	if !ok {
		fmt.Println("error casting to ListPermissionsResultValue")
	}

	r.value = value
}

func NewListPermissionsResult(key uuid.UUID, value ListPermissionsResultValue) *ListPermissionsResult {
	r := new(ListPermissionsResult)
	r.key = key
	r.value = value

	return r
}

type TaskerResultInterface interface {
	GetKey() string
	GetValue() interface{}
	SetValue(interface{})
}

type TaskerJob struct {
	ID uuid.UUID

	RoleID string
	Type   string

	results chan TaskerResultInterface
	wg      *sync.WaitGroup
}

type Tasker struct {
	wg sync.WaitGroup

	// generic config
	workers int

	register sync.Map

	jobs         chan TaskerJob
	shutdownCtx  context.Context
	shutdownFunc context.CancelCauseFunc

	ofga OpenFGAClientInterface
}

func (t *Tasker) Go(roleID, ofgaType string, results chan TaskerResultInterface, wg *sync.WaitGroup) error {
	fmt.Println("submitting task with ", roleID, ofgaType)
	select {
	case t.jobs <- TaskerJob{
		ID:     uuid.New(),
		RoleID: roleID,
		Type:   ofgaType,

		results: results,
		wg:      wg,
	}:
		return nil
	default:
		return fmt.Errorf("error")
	}
}

func (t *Tasker) Shutdown() {
	t.shutdownFunc(fmt.Errorf("shutting down"))
	t.wg.Wait()
}

// TODO @shipperizer move to dynamic pool
func (t *Tasker) start(workers int) {
	// if len(t.register)+workers >= t.workers {
	// 	fmt.Println("max number of workers reached")
	// 	return
	// }

	t.wg.Add(workers)

	for i := 0; i < workers; i++ {
		go t.consume(uuid.New())
	}
}

func (t *Tasker) executeFunc(workerID, taskID, roleID, ofgaType string, results chan TaskerResultInterface, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}

	select {
	case <-t.shutdownCtx.Done():
		fmt.Println("shutdown received")
	default:
		fmt.Println("executing ", taskID, " on worker ", workerID)
		fmt.Println("with roleID ", roleID, " and type ", ofgaType)

		permissions, err := t.listPermissions(context.Background(), roleID, ofgaType)

		results <- NewListPermissionsResult(
			uuid.New(),
			*NewListPermissionsResultValue(permissions, err),
		)
	}
}

func (t *Tasker) listPermissions(ctx context.Context, ID, ofgaType string) ([]string, error) {
	r, err := t.ofga.ReadTuples(ctx, fmt.Sprintf("role:%s#assignee", ID), "", fmt.Sprintf("%s:", ofgaType), "")

	if err != nil {
		return nil, err
	}

	permissions := make([]string, 0)

	for _, t := range r.GetTuples() {
		// if relation doesn't start with can_ it means it's not a permission (see #assignee)
		if !strings.HasPrefix(t.Key.Relation, "can_") {
			continue
		}

		permissions = append(permissions, fmt.Sprintf("%s::%s", t.Key.Relation, t.Key.Object))
	}

	return permissions, nil
}

func (t *Tasker) consume(ID uuid.UUID) {
	t.register.Store(ID, true)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in consume ", ID.String(), " ", r)
			t.register.Delete(ID)
			t.wg.Done()
			return
		}
	}()

	for {
		select {
		case job := <-t.jobs:
			fmt.Println(ID, ": received job: ", job)
			t.executeFunc(ID.String(), job.ID.String(), job.RoleID, job.Type, job.results, job.wg)
			fmt.Println(ID, ": done job: ", job)
		case <-t.shutdownCtx.Done():
			fmt.Println(ID, " going down")
			t.register.Delete(ID)
			t.wg.Done()
			return
		}

	}

}

func NewTasker(workers int) *Tasker {
	t := new(Tasker)

	t.workers = workers

	t.shutdownCtx, t.shutdownFunc = context.WithCancelCause(context.Background())
	t.jobs = make(chan TaskerJob)
	t.ofga = NewOFGAClient()
	go t.start(t.workers)

	return t
}
