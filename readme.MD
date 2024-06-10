# httpmate

httpmate is a command-line tool written in Go that allows you to manage 
collections of HTTP requests. With httpmate, you can create and manage 
collections of HTTP requests, configure individual requests, execute them, and 
inspect their configurations. 

It's a convenient tool for developers who need to manage and test multiple 
HTTP requests efficiently.

## Features

- **Create Collections**: Organize your HTTP requests into collections.
- **Create Requests**: Add HTTP requests to your collections.
- **Edit Request Configurations**: Modify the configuration of existing requests.
- **Make HTTP Requests**: Execute HTTP requests based on the configurations.
- **Remove Requests**: Delete specific requests from your collections.
- **Remove Collections**: Delete entire collections of requests.
- **Inspect Request Configurations**: View the configuration details of a request.
- **List Requests and Collections**: Display all requests and collections.

## Installation

1. Download the latest release of httpmate from the [releases page](https://github.com/joaocgduarte/httpmate/releases).

2. Extract the binary file from the downloaded archive.

3. Move the binary file to a directory that is in your system's `PATH`. For example:

    ```sh
    mv httpmate /usr/local/bin/
    ```

4. Ensure the binary has executable permissions:

    ```sh
    sudo chmod +x /usr/local/bin/httpmate
    ```

## Usage
Upon first usage, it will create configuration files under ~/.config/httpmate,
which are necessary configurations to for httpmate. You can customize them as
you want.

### Create a Request (and collections)
You will be prompted to add all the necessary configurations

```sh
 httpmate create 
```

### Make an HTTP Request
```sh
httpmate run
```

### Remove a Request
```sh
httpmate remove
```

### Remove a collection
```sh
httpmate remove-collection
```

### Inspect a Request Configuration
```sh
httpmate inspect
```

### List requests
```sh
httpmate list
```

You can add --help on any of the commands to get additional information about 
each of the commands.


## Contributing
Contributions are welcome! If you find a bug or want to add a feature, feel 
free to open an issue or a pull request.

1. Fork the repository
2. Create a new branch (git checkout -b feature/your-feature)
3. Commit your changes (git commit -am 'Add new feature')
4. Push to the branch (git push origin feature/your-feature)
5. Open a pull request

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
