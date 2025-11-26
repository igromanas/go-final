package api

import (
	"log"
	"net/http"

	"github.com/igromanas/go-final/pkg/db"
	"github.com/igromanas/go-final/pkg/model"
)

const TaskLimit = 50

type TasksResp struct {
	Tasks []*model.Task `json:"tasks"`
}

/*
В браузере рядом с кнопкой Добавить задачу есть поле для поиска.
Добавьте возможность выбрать задачи через эту строку.
Необходимо проверить наличие строки поиска в заголовке или комментарии задач.
Обработчик должен дополнительно обрабатывать параметр search в строке запроса.
Например, /api/tasks?search=бассейн возвратит задачи со словом «бассейн».
Ещё добавьте возможность выбрать задачи на конкретную дату.
Для этого нужно проверять search на соответствие формату 02.01.2006.
То есть, по запросу /api/tasks?search=08.02.2024 должны возвратиться задачи на 8 февраля 2024 года.
*/

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Printf("[REQ] error: wrong method")
		writeErrorJSON(w, http.StatusMethodNotAllowed, "wrong method error")
		return
	}

	tasks, err := db.Tasks(TaskLimit)
	if err != nil {
		log.Printf("[DB] error getting tasks: %v", err)
		writeErrorJSON(w, http.StatusInternalServerError, "db error")
		return
	}

	writeJSON(w, TasksResp{Tasks: tasks})
}
