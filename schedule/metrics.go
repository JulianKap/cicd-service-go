package schedule

//
//import (
//	"cicd-service-go/manager"
//	"github.com/prometheus/client_golang/prometheus"
//	log "github.com/sirupsen/logrus"
//)
//
//var (
//	//Метрики мастера
//	metricsTotalProjects = prometheus.NewGauge(prometheus.GaugeOpts{
//		Name: "cicd_service_total_projects",
//		Help: "Total number of projects",
//	})
//
//	metricsTotalTasks = prometheus.NewGaugeVec(prometheus.GaugeOpts{
//		Name: "cicd_service_total_tasks",
//		Help: "Total number of issues with project key",
//	}, []string{"project_key"})
//
//	metricsTotalJobs = prometheus.NewGaugeVec(prometheus.GaugeOpts{
//		Name: "cicd_service_total_jobs",
//		Help: "Total number of tasks with project key and tasks by status",
//	}, []string{"project_key", "task_key", "status"})
//
//	metricsRunningJobs = prometheus.NewGauge(prometheus.GaugeOpts{
//		Name: "cicd_service_running_jobs",
//		Help: "Number of running jobs",
//	})
//
//	metricsFailedJobs = prometheus.NewGauge(prometheus.GaugeOpts{
//		Name: "cicd_service_failed_jobs",
//		Help: "Number of erroneous jobs",
//	})
//)
//
//// Регистрация метрик в зависимости от роли
//func registerMetrics() {
//	prometheus
//
//	// Регистрация метрик
//	prometheus.MustRegister(metricsTotalProjects)
//	prometheus.MustRegister(metricsTotalTasks)
//	prometheus.MustRegister(metricsTotalJobs)
//	prometheus.MustRegister(metricsRunningJobs)
//	prometheus.MustRegister(metricsFailedJobs)
//
//	if manager.MemberInfo.Role == manager.MasterRole || manager.MemberInfo.Standalone {
//		prometheus.MustRegister(totalProjects)
//		prometheus.MustRegister(totalTasks)
//		prometheus.MustRegister(totalJobs)
//	} else if manager.MemberInfo.Role == manager.WorkerRole {
//		prometheus.MustRegister(runningJobs)
//		prometheus.MustRegister(failedJobs)
//	}
//}
//
//// Функция для обновления метрик в зависимости от роли
//func updateMetrics(projects []Project) {
//
//	if manager.MemberInfo.Role == manager.MasterRole && !manager.MemberInfo.Standalone {
//		log.Debug("runScheduleWorker #1: in master rule")
//		return nil
//	}
//
//	if role == "master" {
//		totalProjects.Set(float64(len(projects)))
//
//		totalTasks.Reset()
//		totalJobs.Reset()
//
//		for _, project := range projects {
//			totalTasks.WithLabelValues(project.Key).Set(float64(len(project.Tasks)))
//
//			for _, task := range project.Tasks {
//				for _, job := range task.Jobs {
//					totalJobs.WithLabelValues(project.Key, task.Key, job.Status).Inc()
//				}
//			}
//		}
//	} else if role == "worker" {
//		runningJobs.Set(0)
//		failedJobs.Set(0)
//
//		for _, project := range projects {
//			for _, task := range project.Tasks {
//				for _, job := range task.Jobs {
//					if job.Status == "running" {
//						runningJobs.Inc()
//					} else if job.Status == "failed" {
//						failedJobs.Inc()
//					}
//				}
//			}
//		}
//	}
//}
