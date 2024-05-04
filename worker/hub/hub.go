package hub

//
//var workerHubInstance workerHub
//
//func init() {
//	workerHubInstance = workerHub{
//		controlChan: make(chan int, 1),
//	}
//
//	//handlers.NewWorker(&workerHubInstance)
//}
//
//func GetHub() worker.IWorkerHub {
//	return &workerHubInstance
//}
//
//func (wh *workerHub) GetWorker(workerName string) worker.IWorker {
//	switch workerName {
//	case handlers.CWorkerType:
//		return handlers.GetWorker()
//	}
//
//	return nil
//}
//
//func (wh *workerHub) ChanLock() {
//	wh.controlChan <- 1
//}
//
//func (wh *workerHub) ChanUnlock() {
//	<-wh.controlChan
//}
//
//func (wh *workerHub) GetStatusRunningWork() (taskpkg.ITask, bool, error) {
//	if wh.runningWorker == nil {
//		return nil, false, nil
//	}
//
//	var (
//		task      taskpkg.ITask
//		isRunning bool
//		err       error
//	)
//
//	defer func() {
//		if !isRunning {
//			wh.DelRunningWork()
//		}
//	}()
//
//	task, isRunning, err = wh.runningWorker.StatusWork()
//	if err != nil {
//		return nil, false, errors.Wrap(err, "get status work from running worker")
//	}
//
//	return task, isRunning, nil
//}
//
//func (wh *workerHub) SetRunningWorker(worker worker.IWorker) {
//	wh.runningWorker = worker
//}
//
//func (wh *workerHub) GetRunningWorker() worker.IWorker {
//	return wh.runningWorker
//}
//
//func (wh *workerHub) DelRunningWork() {
//	wh.runningWorker = nil
//}
