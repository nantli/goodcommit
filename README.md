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

## Creating Your Own Commiter

If using the default commiter is not good enough for you, you can create your own commiter that implements the `Commiter` interface. This allows you to define how the commit form operates, handles input, or processes the final commit message. Here's how:

1. **Define Your Commiter**: Create a new Go file in `pkg/commiters/yourcommiter/`. Define a struct that implements the `Commiter` interface from `pkg/commiter/commiter.go`.

   ```go
   package yourcommiter

   import "github.com/nantli/goodcommit/pkg/commiter"

   type YourCommiter struct {
       // Your fields here
   }

   func (yc *YourCommiter) RunForm(accessible bool) error {
       // Implement how your commiter runs the form
   }

   func (yc *YourCommiter) RunPostProcessing() error {
       // Implement any post-processing steps
   }

   func (yc *YourCommiter) PreviewCommit() {
       // Implement how to preview the commit
   }

   func (yc *YourCommiter) Commit() error {
       // Implement the commit action
   }
   ```

2. **Register Your Commiter**: In your `main.go`, import your commiter and use it when creating a new `GoodCommit` instance.

   ```go:cmd/main.go
   import (
       "github.com/nantli/goodcommit/pkg/goodcommit"
       "github.com/nantli/goodcommit/pkg/commiters/yourcommiter"
   )

   func main() {
       // Initialize your commiter
       yourCommiter := yourcommiter.New()
       // Create a GoodCommit instance with your commiter
       goodCommit := goodcommit.New(yourCommiter)
       // Execute
       if err := goodCommit.Execute(accessible); err != nil {
           fmt.Println("Error occurred:", err)
           os.Exit(1)
       }
   }
   ```

3. **Test Your Commiter**: After implementing your commiter, test it thoroughly to ensure it works as expected with the `goodcommit` form.

By following these steps, you can create and integrate your own commiter into `goodcommit`, allowing for a highly customized commit process.
