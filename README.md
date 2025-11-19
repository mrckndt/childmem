# childmem

`childmem` is a command-line tool that monitors the memory usage of child processes of a specified parent process on Linux systems using `procfs`.

## Building
```bash
go build -o childmem .
```

## Usage

The `childmem` tool requires either a parent process name (`-pname`) or a parent process ID (`-ppid`) to identify the target process. It outputs the child process information to a CSV file, defaulting to `child_mem.csv`.

```bash
./childmem -pname <parent_process_name> [-output <output_file.csv>]
./childmem -ppid <parent_process_id> [-output <output_file.csv>]
```

### Arguments

-   `-pname <name>`: The name of the parent process to monitor.
-   `-ppid <pid>`: The PID of the parent process to monitor.
-   `-includeParent`: (Optional) Include the parent process in the output.
-   `-output <file>`: (Optional) Path to the output CSV file. Defaults to `./child_mem.csv`.

## Docker Usage
### Building the image
```bash
docker build -t childmem -f build/Dockerfile .
```

### Run with defaults (30s interval, monitoring "mattermost" process)
```bash
docker run --security-opt=no-new-privileges --pid=host -v ./data:/app/data childmem
```

### Run with custom settings
```bash
# Custom interval and process name
docker run -e INTERVAL=10 -e PNAME="nginx" --security-opt=no-new-privileges --pid=host -v ./data:/app/data childmem
```

### Environment variables
- `INTERVAL`: Seconds between runs (default: 30)
- `PNAME`: Process name to monitor (default: "mattermost")

### Security Notes
- `childmem` runs as non-root user (UID/GID 1000) inside container
- Requires `--pid=host` to access host process information
- Ensure the data directory on the host is owned by UID 1000 (`sudo chown -R 1000:1000 ./data`)
