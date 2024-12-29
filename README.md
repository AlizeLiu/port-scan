# port-scan

`port-scan` 是一个用 Go 语言编写的命令行工具，用于扫描指定 IP 地址上的端口，并显示开放端口的信息。

## 功能概述

- 扫描指定范围内的端口，识别哪些端口是开放的。
- 显示端口的运行状态（开放或关闭）。
- 通过调用 `lsof` 工具尝试识别端口运行的服务。
- 支持多线程扫描以提高效率。

## 安装

1. 确保您的系统安装了 Go 编译器（版本 >= 1.18）。
2. 克隆或下载项目代码。
3. 在项目根目录下运行以下命令构建可执行文件：

   ```bash
   go build -o port-scan
   ```

## 使用方法

### 基本命令

运行以下命令查看工具的帮助信息：

```bash
./port-scan --help
```

输出示例：

```
port-scan is a tool for scanning ports

Usage:
  port-scan [command]

Available Commands:
  scan        Scan ports on a given IP address

Flags:
  -h, --help   help for port-scan
```

### 端口扫描

使用 `scan` 子命令执行端口扫描。示例如下：

```bash
./port-scan scan -i <IP地址> -s <起始端口> -e <结束端口>
```

#### 参数说明

- `-i, --ip`：指定要扫描的目标 IP 地址（必选）。
- `-s, --start-port`：指定扫描的起始端口，默认为 `1`。
- `-e, --end-port`：指定扫描的结束端口，默认为 `65535`。

#### 示例

扫描 `192.168.1.1` 的 20 到 80 端口：

```bash
./port-scan scan -i 192.168.1.1 -s 20 -e 80
```

输出示例：

```
Scanning IP: 192.168.1.1 from port 20 to port 80
端口	状态	运行的服务
----------------------------
22	open	ssh
80	open	http
```

### 扫描结果说明

- **端口**：扫描到的端口号。
- **状态**：端口状态，`open` 表示开放，`closed` 表示关闭。
- **运行的服务**：如果工具能够识别到服务，将显示服务名称，否则显示 `unknown`。

## 实现细节

1. **端口扫描**：通过 `net.DialTimeout` 尝试连接目标端口，判断其是否开放。
2. **服务识别**：调用系统工具 `lsof` 并解析其输出，尝试识别开放端口对应的服务。
3. **多线程**：使用 Go 的 Goroutine 和 Channel 实现高效并发扫描。

## 依赖

- [cobra](https://github.com/spf13/cobra)：用于处理命令行参数。

可以通过以下命令安装 `cobra`：

```bash
go get -u github.com/spf13/cobra
```

## 注意事项

- 工具需要权限访问目标端口和运行 `lsof`，建议使用具有管理员权限的用户执行。
- `lsof` 命令可能因系统环境差异产生不同的输出格式，请确保在兼容的环境中运行。

## 许可证

该项目使用 MIT 许可证。详情参见 LICENSE 文件。
