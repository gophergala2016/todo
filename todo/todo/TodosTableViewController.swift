
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
        self.navigationItem.leftBarButtonItem = UIBarButtonItem(barButtonSystemItem: .Refresh, target: self, action: "refresh")
        self.navigationItem.rightBarButtonItem = UIBarButtonItem(barButtonSystemItem: .Add, target: self, action: "add")
        
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
        cell.textLabel?.text = todos[indexPath.row].text()
        cell.accessoryType = UITableViewCellAccessoryType.Checkmark
        return cell
    }
    
    func getTodos() {
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
    
    func add() {
        let vc = AddTodoViewController()
        let nav = UINavigationController(rootViewController: vc)
        self.presentViewController(nav, animated: true, completion: nil)
    }
    
    func refresh() {
        let manager = CBLManager.sharedInstance()
        do {
            let database = try manager.databaseNamed("todos")
            let url = NSURL(string: "http://192.168.99.100:5984/john/")
            self.push = database.createPushReplication(url!)
            self.pull = database.createPullReplication(url!)
            NSNotificationCenter.defaultCenter().addObserver(self, selector: "replicationChanged:", name: kCBLReplicationChangeNotification, object: push)
            NSNotificationCenter.defaultCenter().addObserver(self, selector: "replicationChanged:", name: kCBLReplicationChangeNotification, object: pull)
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
                if push.status == .Stopped && pull.status == .Stopped {
                    print("done")
                    print(pull.lastError)
                    print(push.lastError)
                    self.getTodos()
                }
        }
    }

}
