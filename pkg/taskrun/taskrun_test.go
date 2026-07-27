// Copyright © 2019 The Tekton Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package taskrun

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSortTasksBySpecOrder_Matrix(t *testing.T) {
	now := metav1.Now()

	pipelineTasks := []v1.PipelineTask{
		{Name: "matrixed-task"},
		{Name: "regular-task"},
	}

	taskRuns := map[string]*v1.PipelineRunTaskRunStatus{
		"matrixed-task-0": {
			PipelineTaskName: "matrixed-task",
			Status: &v1.TaskRunStatus{
				TaskRunStatusFields: v1.TaskRunStatusFields{StartTime: &now},
			},
		},
		"matrixed-task-1": {
			PipelineTaskName: "matrixed-task",
			Status: &v1.TaskRunStatus{
				TaskRunStatusFields: v1.TaskRunStatusFields{StartTime: &now},
			},
		},
		"regular-task": {
			PipelineTaskName: "regular-task",
			Status: &v1.TaskRunStatus{
				TaskRunStatusFields: v1.TaskRunStatusFields{StartTime: &now},
			},
		},
	}

	got := SortTasksBySpecOrder(pipelineTasks, taskRuns)

	wantNames := []string{"matrixed-task-0", "matrixed-task-1", "regular-task"}
	gotNames := []string{}
	for _, r := range got {
		gotNames = append(gotNames, r.Name)
	}

	if d := cmp.Diff(wantNames, gotNames, cmpopts.SortSlices(func(a, b string) bool { return a < b })); d != "" {
		t.Fatalf("SortTasksBySpecOrder() did not return all TaskRuns fanned out from a matrixed PipelineTask: %s", d)
	}

	for _, r := range got {
		if r.Name == "matrixed-task-0" || r.Name == "matrixed-task-1" {
			if r.Task != "matrixed-task" {
				t.Fatalf("expected Run %q to belong to PipelineTask %q, got %q", r.Name, "matrixed-task", r.Task)
			}
		}
	}
}
