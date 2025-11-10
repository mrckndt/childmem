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
