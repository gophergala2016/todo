
import UIKit
import Item

class TodosTableViewController: UITableViewController {
    
    private var pull: CBLReplication?
    private var push: CBLReplication?
    private var todos: [GoItemTodo] = []
    private let reuseIdentifier = "reuseIdentifier"

    override func viewDidLoad() {
        super.viewDidLoad()
        self.title = "Todos"
        self.tableView.registerClass(UITableViewCell.self, forCellReuseIdentifier: reuseIdentifier)
        self.navigationItem.rightBarButtonItem = UIBarButtonItem(barButtonSystemItem: .Add, target: self, action: "add")
        
        // start syncing
        self.sync()
        
        // load todos
        self.getTodos()
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
        let todo = todos[indexPath.row]
        cell.textLabel?.text = todo.text()
        if todo.done() {
            cell.accessoryType = .Checkmark
        } else {
            cell.accessoryType = .None
        }
        return cell
    }
    
    override func tableView(tableView: UITableView, editActionsForRowAtIndexPath indexPath: NSIndexPath) -> [UITableViewRowAction]? {
        let done = UITableViewRowAction(style: .Normal, title: "Done", handler: { action, indexPath in
            do {
                let database = try Database.Get()
                let todo = Database.GetTodoByID(database, id: self.todos[indexPath.row].id())
                todo.setDone(true)
                try Database.UpdateTodo(database, item: todo)
                self.tableView.setEditing(false, animated: true)
                self.getTodos()
            } catch {
                print(error)
            }
        })
        let undone = UITableViewRowAction(style: .Normal, title: "Undone", handler: { action, indexPath in
            do {
                let database = try Database.Get()
                let todo = Database.GetTodoByID(database, id: self.todos[indexPath.row].id())
                todo.setDone(false)
                try Database.UpdateTodo(database, item: todo)
                self.tableView.setEditing(false, animated: true)
                self.getTodos()
            } catch {
                print(error)
            }
        })
        let delete = UITableViewRowAction(style: .Default, title: "Remove", handler: { action, indexPath in
            do {
                let database = try Database.Get()
                try Database.RemoveTodoByID(database, id: self.todos[indexPath.row].id())
                self.tableView.setEditing(false, animated: true)
                self.getTodos()
            } catch {
                print(error)
            }
        })
        undone.backgroundColor = UIColor.blueColor()
        done.backgroundColor = UIColor.blueColor()
        delete.backgroundColor = UIColor.orangeColor()
        if todos[indexPath.row].done() {
            return [delete, undone]
        } else {
            return [delete, done]
        }
    }
    
    func getTodos() {
        // todos from database
        do {
            let database = try Database.Get()
            database.manager.backgroundTellDatabaseNamed(database.name, to: { bgdb in
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
    
    func add() {
        let vc = AddTodoViewController()
        let nav = UINavigationController(rootViewController: vc)
        self.presentViewController(nav, animated: true, completion: nil)
    }
    
    func sync() {
        do {
            let database = try Database.Get()
            let url = NSURL(string: "http://\(Config.URL):5984/john/")
            self.push = database.createPushReplication(url!)
            self.pull = database.createPullReplication(url!)
            NSNotificationCenter.defaultCenter().addObserver(self, selector: "replicationChanged:", name: kCBLReplicationChangeNotification, object: push)
            NSNotificationCenter.defaultCenter().addObserver(self, selector: "replicationChanged:", name: kCBLReplicationChangeNotification, object: pull)
            self.push?.continuous = true
            self.pull?.continuous = true
            self.push!.start()
            self.pull!.start()
        } catch {
            print(error)
        }
    }
    
    func replicationChanged(n: NSNotification) {
        if
            let pull = self.pull,
            let push = self.push {
                if push.status == .Idle && pull.status == .Idle {
                    print(pull.lastError)
                    print(push.lastError)
                    self.getTodos()
                }
        }
    }

}
