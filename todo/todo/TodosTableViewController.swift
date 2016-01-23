
import UIKit
import Item

class TodosTableViewController: UITableViewController {
    
    private var todos: [GoItemTodo] = []
    private let reuseIdentifier = "reuseIdentifier"

    override func viewDidLoad() {
        super.viewDidLoad()
        
        self.title = "Todos"
        self.tableView.registerClass(UITableViewCell.self, forCellReuseIdentifier: reuseIdentifier)
        self.navigationItem.rightBarButtonItem = UIBarButtonItem(barButtonSystemItem: .Add, target: self, action: "add")
        
        // todos from database
        let manager = CBLManager.sharedInstance()
        do {
            let database = try manager.databaseNamed("todos")
            manager.backgroundTellDatabaseNamed(database.name, to: { bgdb in
                do {
                    let todos = try Database.GetTodos(bgdb)
                    dispatch_async(dispatch_get_main_queue(), {
                        self.todos = todos
                        self.tableView.reloadData()
                    })
                } catch {
                    print(error)
                }
            })
        } catch {
            print(error)
        }
    }

    override func didReceiveMemoryWarning() {
        super.didReceiveMemoryWarning()
    }

    override func numberOfSectionsInTableView(tableView: UITableView) -> Int {
        return 1
    }

    override func tableView(tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
        return todos.count
    }

    override func tableView(tableView: UITableView, cellForRowAtIndexPath indexPath: NSIndexPath) -> UITableViewCell {
        let cell = tableView.dequeueReusableCellWithIdentifier(reuseIdentifier, forIndexPath: indexPath)
        cell.textLabel?.text = todos[indexPath.row].text()
        return cell
    }
    
    func add() {
        let vc = AddTodoViewController()
        let nav = UINavigationController(rootViewController: vc)
        self.presentViewController(nav, animated: true, completion: nil)
    }

}
