package redis

import (
	"strings"

	"github.com/goravel/framework/contracts/http"

	"panel/app/http/controllers"
	"panel/pkg/tools"
)

type RedisController struct {
}

type Info struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func NewRedisController() *RedisController {
	return &RedisController{}
}

// Status 获取运行状态
func (c *RedisController) Status(ctx http.Context) {
	if !controllers.Check(ctx, "redis") {
		return
	}

	status := tools.ExecShell("systemctl status redis | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取Redis状态失败")
		return
	}

	if status == "active" {
		controllers.Success(ctx, true)
	} else {
		controllers.Success(ctx, false)
	}
}

// Reload 重载配置
func (c *RedisController) Reload(ctx http.Context) {
	if !controllers.Check(ctx, "redis") {
		return
	}

	tools.ExecShell("systemctl reload redis")
	status := tools.ExecShell("systemctl status redis | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取Redis状态失败")
		return
	}

	if status == "active" {
		controllers.Success(ctx, true)
	} else {
		controllers.Success(ctx, false)
	}
}

// Restart 重启服务
func (c *RedisController) Restart(ctx http.Context) {
	if !controllers.Check(ctx, "redis") {
		return
	}

	tools.ExecShell("systemctl restart redis")
	status := tools.ExecShell("systemctl status redis | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取Redis状态失败")
		return
	}

	if status == "active" {
		controllers.Success(ctx, true)
	} else {
		controllers.Success(ctx, false)
	}
}

// Start 启动服务
func (c *RedisController) Start(ctx http.Context) {
	if !controllers.Check(ctx, "redis") {
		return
	}

	tools.ExecShell("systemctl start redis")
	status := tools.ExecShell("systemctl status redis | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取Redis状态失败")
		return
	}

	if status == "active" {
		controllers.Success(ctx, true)
	} else {
		controllers.Success(ctx, false)
	}
}

// Stop 停止服务
func (c *RedisController) Stop(ctx http.Context) {
	if !controllers.Check(ctx, "redis") {
		return
	}

	tools.ExecShell("systemctl stop redis")
	status := tools.ExecShell("systemctl status redis | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取Redis状态失败")
		return
	}

	if status != "active" {
		controllers.Success(ctx, true)
	} else {
		controllers.Success(ctx, false)
	}
}

// GetConfig 获取配置
func (c *RedisController) GetConfig(ctx http.Context) {
	if !controllers.Check(ctx, "redis") {
		return
	}

	// 获取配置
	config := tools.ReadFile("/www/server/redis/redis.conf")
	if len(config) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取Redis配置失败")
		return
	}

	controllers.Success(ctx, config)
}

// SaveConfig 保存配置
func (c *RedisController) SaveConfig(ctx http.Context) {
	if !controllers.Check(ctx, "redis") {
		return
	}

	config := ctx.Request().Input("config")
	if len(config) == 0 {
		controllers.Error(ctx, http.StatusBadRequest, "配置不能为空")
		return
	}

	if !tools.WriteFile("/www/server/redis/redis.conf", config, 0644) {
		controllers.Error(ctx, http.StatusInternalServerError, "写入Redis配置失败")
		return
	}

	tools.ExecShell("systemctl restart redis")

	controllers.Success(ctx, nil)
}

// Load 获取负载
func (c *RedisController) Load(ctx http.Context) {
	if !controllers.Check(ctx, "redis") {
		return
	}

	status := tools.ExecShell("systemctl status redis | grep Active | grep -v grep | awk '{print $2}'")
	if status != "active" {
		controllers.Error(ctx, http.StatusInternalServerError, "Redis 已停止运行")
		return
	}

	raw := tools.ExecShell("redis-cli info")
	if len(raw) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取Redis负载失败")
		return
	}

	infoLines := strings.Split(raw, "\n")
	dataRaw := make(map[string]string)

	for _, item := range infoLines {
		parts := strings.Split(item, ":")
		if len(parts) == 2 {
			dataRaw[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	data := []Info{
		{"TCP 端口", dataRaw["tcp_port"]},
		{"已运行天数", dataRaw["uptime_in_days"]},
		{"连接的客户端数", dataRaw["connected_clients"]},
		{"已分配的内存总量", dataRaw["used_memory_human"]},
		{"占用内存总量", dataRaw["used_memory_rss_human"]},
		{"占用内存峰值", dataRaw["used_memory_peak_human"]},
		{"内存碎片比率", dataRaw["mem_fragmentation_ratio"]},
		{"运行以来连接过的客户端的总数", dataRaw["total_connections_received"]},
		{"运行以来执行过的命令的总数", dataRaw["total_commands_processed"]},
		{"每秒执行的命令数", dataRaw["instantaneous_ops_per_sec"]},
		{"查找数据库键成功次数", dataRaw["keyspace_hits"]},
		{"查找数据库键失败次数", dataRaw["keyspace_misses"]},
		{"最近一次 fork() 操作耗费的毫秒数", dataRaw["latest_fork_usec"]},
	}

	controllers.Success(ctx, data)
}
