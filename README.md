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
   ./goodcommit
   ```

## Developing New Modules

Modules in `goodcommit` allow for extensibility and customization of the commit form. To develop a new module, follow these steps:

1. **Create a New Module File**: In `pkg/modules/`, create a new Go file for your module, e.g., `mymodule.go`.

2. **Define Your Module Struct**: Implement the `module.Module` interface.
   ```go
   package mymodule

   import (
       "github.com/charmbracelet/huh"
       "github.com/nantli/goodcommit/pkg/module"
   )

   type MyModule struct {
       config module.Config
   }

   // Implement interface methods: LoadConfig, NewField, PostProcess, etc.
   ```

3. **Implement Required Methods**: At minimum, implement `LoadConfig`, `NewField`, and `PostProcess` methods as per your module's functionality.

4. **Register Your Module**: In `cmd/main.go`, import your module and add it to the `modules` slice.
   ```go
   import (
       "github.com/nantli/goodcommit/pkg/modules/mymodule"
   )

   func main() {
       // Other modules...
       myModule := mymodule.New()
       modules = append(modules, myModule)
       // Continue setup...
   }
   ```

## Creating Your Own Commit Formatting

To customize `goodcommit` further, you might want to modify the `commiter` to suit your specific needs. Here's how:

1. **Modify Commiter**: In `pkg/commiter/commiter.go`, you can add or modify methods to change how the commit form operates, handles input, or processes the final commit message.

2. **Adjust Module Interaction**: If your customization involves new ways modules interact with each other or with the `commiter`, ensure to update those interactions in the respective module files.

3. **Test Your Changes**: After making modifications, test `goodcommit` thoroughly to ensure your changes work as expected.
