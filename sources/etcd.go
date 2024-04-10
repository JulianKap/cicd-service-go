package sources

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strconv"
	"time"
)

// Добавляем project
func (p *Project) createProject(cli *clientv3.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Преобразование структуры Project в JSON
	projectJSON, err := json.Marshal(p)
	if err != nil {
		log.Error("createProject[0]: Ошибка преобразования Project в json: ", err)
		return err
	}

	// Добавление данных в etcd
	_, err = cli.Put(ctx, "/service-config/projects/"+p.ID, string(projectJSON))
	if err != nil {
		log.Error("createProject[1]: Ошибка добавления структуры Project в etcd: ", err)
		return err
	}

	return nil
}

// Добавляем job
func (j *Job) createJob(cli *clientv3.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Преобразование структуры Job в JSON
	jobJSON, err := json.Marshal(j)
	if err != nil {
		log.Error("createJob[0]: Ошибка преобразования Job в json: ", err)
		return err
	}

	// Добавление данных в etcd
	_, err = cli.Put(ctx, fmt.Sprintf("/service-config/projects/%s/jobs/%s", j.IdProject, j.ID), string(jobJSON))
	if err != nil {
		log.Error("createJob[1]: Ошибка добавления структуры Job в etcd: ", err)
		return err
	}

	return nil
}

// getLastID возвращает последний ID из etcd или 1, если ничего не найдено
func getLastID(cli *clientv3.Client, key string) (int64, error) {
	// Контекст для операций
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Получение значения ключа из etcd
	resp, err := cli.Get(ctx, key)
	if err != nil {
		return 1, nil // В случае ошибки возвращаем 1
	}

	// Если ключ существует и содержит значение, попытаемся преобразовать его в число
	if len(resp.Kvs) > 0 {
		value := string(resp.Kvs[0].Value)
		lastID, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return 1, fmt.Errorf("ошибка преобразования значения ключа в число: %v", err)
		}
		// Возвращаем последний ID, увеличенный на 1
		return lastID + 1, nil
	}

	// Если ключ не найден, возвращаем 1
	return 1, nil
}
