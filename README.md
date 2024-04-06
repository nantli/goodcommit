# You're gonna like this commit, he's all right. He's a commit

Welcome to `goodcommit` - a customizable commit message generator that ensures your commit messages follow best practices and are consistent across your projects.

## Building the Project

To build `goodcommit`, ensure you have Go installed on your system. Then, follow these steps:

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/goodcommit.git
   ```
2. Navigate to the project directory:
   ```bash
   cd goodcommit
   ```
3. Build the project:
   ```bash
   go build -o goodcommit ./cmd/main.go
   ```
4. Run `goodcommit`:
   ```bash
   git add .
   ./goodcommit
   ```

### Specifying a Configuration File

To use a custom configuration file with `goodcommit`, you have two options:

1. **Using a Command Line Flag**: Specify the `--config` flag followed by the path to your configuration file when running the program:
   ```bash
   ./goodcommit --config /path/to/your/config.json
   ```

2. **Using an Environment Variable**: Set the `GOODCOMMIT_CONFIG_PATH` environment variable to the path of your configuration file. If both the environment variable and the `--config` flag are provided, the `--config` flag takes precedence.
   ```bash
   export GOODCOMMIT_CONFIG_PATH=/path/to/your/config.json
   ./goodcommit
   ```

## Developing New Modules

Modules in `goodcommit` allow for extensibility and customization of the commit form. To develop a new module, follow these steps:

1. **Create a New Module File**: In the root directory of the project, create a new Go file for your module, e.g., `mymodule.go`.

2. **Define Your Module Struct**: Implement the `gc.Module` interface.

```go
package mymodule

import (
    "github.com/charmbracelet/huh"
    gc "github.com/nantli/goodcommit"
)

type myModule struct {
    config gc.ModuleConfig
}

// Implement interface methods: LoadConfig, NewField, PostProcess, etc.
```


3. **Implement Required Methods**: At minimum, implement `LoadConfig`, `NewField`, `PostProcess`, `Config`, `Name`, `InitCommitInfo`, and `IsActive` methods as per your module's functionality.

4. **Register Your Module**: In your own implementation of `cmd/goodcommit/main.go`, import your goodcommit module and add it to the `modules` slice.

```go
import (
    "github.com/nantli/goodcommit/mymodule"
)

func main() {

// Other modules...
myModule := mymodule.New()
modules = append(modules, myModule)

// Continue setup...
```


## Creating Your Own Commiter

If using the default commiter is not good enough for you, you can create your own commiter that implements the `gc.Commiter` interface. This allows you to define how the commit form operates, handles input, or processes the final commit message. Here's how:

1. **Define Your Commiter**: Create a new Go file in the root directory of the project, e.g., `yourcommiter.go`. Define a struct that implements the `gc.Commiter` interface.

```go
package yourcommiter

import gc "github.com/nantli/goodcommit"

type yourCommiter struct {
    // Your fields here
}

func (yc yourCommiter) LoadModules(modules []gc.Module) error {
    // Implement how your commiter loads the modules
}

func (yc yourCommiter) RunForm(accessible bool) error {
    // Implement how your commiter runs the form
}

func (yc yourCommiter) RunPostProcessing() error {
    // Implement any post-processing steps 
}

func (yc yourCommiter) PreviewCommit() {
    // Implement how your commiter previews the commit
}

func (yc yourCommiter) RenderMessage() string {
    // Implement how your commiter renders the final commit message
}
```

2. **Use Your Commiter**: In your own implementation of `cmd/goodcommit/main.go`, replace the default commiter with your custom commiter.

```go:cmd/goodcommit/main.go
import (
    gc "github.com/nantli/goodcommit"
    "github.com/username/module/yourcommiter"
)

func main() {
    // ...
    // Initialize your commiter
    yourCommiter := yourcommiter.New()
    err = yourCommiter.LoadModules(modules)
    if err != nil {
        fmt.Println("Error occurred while loading modules:", err) os.Exit(1) 
    }

    goodcommit := gc.New(yourCommiter)
    // Execute
    if err := goodcommit.Execute(accessible); err != nil {
        fmt.Println("Error occurred:", err)
        os.Exit(1)
    }
}
```

3. **Test Your Commiter**: After implementing your commiter, test it thoroughly to ensure it works as expected with the `goodcommit` form.

By following these steps, you can create and integrate your own commiter into `goodcommit`, allowing for a highly customized commit process.

## Configuring Modules

Modules in `goodcommit` can be customized through a JSON file. This section outlines how to configure modules, detailing the available fields, their types, and functionalities.

### Configuration Fields

Each module configuration can include the following fields:

- `name`: `string` - The unique identifier for the module.
- `page`: `int` - Determines on which page the module appears in the form.
- `position`: `int` (optional, default: `0`) - The order of the module on the page.
- `pinned`: `bool` (optional, default: `false`) - If `true`, the module is pinned to the top of every page after its initial appearance.
- `active`: `bool` (optional, default: `true`) - Controls the module's activation state. Inactive modules are not displayed.
- `path`: `string` (optional) - Specifies a path to additional configuration or data files required by the module.
- `priority`: `int` (optional, default: `0`) - Used to determine the module's priority. Lower values indicate higher priority.
- `checkpoint`: `bool` (optional, default: `false`) - If `true`, the form will prompt for confirmation before proceeding past this module.
- `dependencies`: `[]string` (optional) - A list of module names that must be active for this module to be activated. This ensures that the current module's functionality is only available if its dependencies are met.

### Examples

Below are examples of different module configurations and their effects:

1. **Basic Module Configuration**

   ```json
   {
     "name": "description",
     "page": 1,
     "position": 1,
     "pinned": false,
     "active": true
   }
   ```
   This configuration activates the `description` module, placing it first on page 1 without pinning it.

2. **Pinned Module Configuration**

   ```json
   {
     "name": "logo",
     "page": 1,
     "position": 1,
     "pinned": true,
     "active": true
   }
   ```
   The `logo` module is activated, pinned, and placed at the top of every page.

3. **Module with External Configuration**

   ```json
   {
     "name": "types",
     "page": 1,
     "position": 2,
     "active": true,
     "path": "./configs/commit_types.json"
   }
   ```
   Activates the `types` module, using an external file for additional configuration.

4. **Module with Dependencies**

   ```json
   {
     "name": "breakingmsg",
     "page": 4,
     "position": 1,
     "active": true,
     "dependencies": ["breaking"]
   }
   ```
   This configuration activates the `breakingmsg` module, which depends on the `breaking` module being active. If the `breaking` module is not active, `breakingmsg` will not be activated.

By adjusting these fields in the `config.json` file, you can tailor the `goodcommit` form to meet your project's specific needs.

### Example Configuration File

This is an example configuration file that activates the modules: `types`, `scopes`, `description`, `body`, `breaking` and `breakingmsg`. You can use this as a starting point for your own configuration.

```json
{
    "activeModules": [
        {
            "name": "types",
            "page": 1,
            "position": 1,
            "active": true,
            "path": "./configs/commit_types.example.json",
            "checkpoint": true
        },
        {
            "name": "scopes",
            "page": 2,
            "position": 1,
            "active": true,
            "path": "./configs/commit_scopes.example.json",
            "dependencies": ["types"],
            "priority": 3
        },
        {
            "name": "description",
            "page": 3,
            "position": 1,
            "active": true
        },
        {
            "name": "body",
            "page": 3,
            "position": 2,
            "active": true,
            "priority": 2
        },
        {
            "name": "breaking",
            "page": 3,
            "position": 3,
            "active": true,
            "priority": 4,
            "checkpoint": true
        },
        {
            "name": "breakingmsg",
            "page": 4,
            "position": 1,
            "active": true,
            "priority": 5,
            "dependencies": ["breaking"]
        }
    ]
}

```
