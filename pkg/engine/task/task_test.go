package task

import (
	"reflect"
	"testing"

	"github.com/kudobuilder/kudo/pkg/apis/kudo/v1alpha1"
	"sigs.k8s.io/yaml"
)

func TestBuild(t *testing.T) {
	type args struct {
		task *v1alpha1.Task
	}
	tests := []struct {
		name     string
		taskYaml string
		want     Tasker
		wantErr  bool
	}{
		{
			name: "apply task",
			taskYaml: `
name: apply-task
kind: Apply
spec: 
    applyResources:
      - pod.yaml
      - service.yaml`,
			want: ApplyTask{
				Name:      "apply-task",
				Resources: []string{"pod.yaml", "service.yaml"},
			},
			wantErr: false,
		},
		{
			name: "delete task",
			taskYaml: `
name: delete-task
kind: Delete
spec: 
    deleteResources:
      - pod.yaml
      - service.yaml`,
			want: DeleteTask{
				Name:      "delete-task",
				Resources: []string{"pod.yaml", "service.yaml"},
			},
			wantErr: false,
		},
		{
			name: "dummy task",
			taskYaml: `
name: dummy-task
kind: Dummy
spec: 
    wantErr: true`,
			want: DummyTask{
				WantErr: true,
			},
			wantErr: false,
		},
		{
			name: "unknown task",
			taskYaml: `
name: unknown-task
kind: Unknown
spec: 
    known: false`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &v1alpha1.Task{}
			err := yaml.Unmarshal([]byte(tt.taskYaml), task)
			if err != nil {
				t.Errorf("Failed to unmarshal task yaml %s: %v", tt.taskYaml, err)
			}

			got, err := Build(task)
			if (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Build() got = %v, want %v", got, tt.want)
			}
		})
	}
}