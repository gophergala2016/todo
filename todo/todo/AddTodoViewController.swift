
import UIKit
import Item

class AddTodoViewController: UIViewController {
    
    private var textfield: UITextField!

    override func viewDidLoad() {
        super.viewDidLoad()

        self.title = "Add todo"
        self.view.backgroundColor = UIColor.whiteColor()
        self.navigationItem.leftBarButtonItem = UIBarButtonItem(barButtonSystemItem: .Cancel, target: self, action: "cancel")
        self.navigationItem.rightBarButtonItem = UIBarButtonItem(barButtonSystemItem: .Done, target: self, action: "done")
        
        // add textfield
        textfield = UITextField()
        textfield.translatesAutoresizingMaskIntoConstraints = false
        textfield.autocorrectionType = .No
        textfield.placeholder = "take the dog for a walk"
        self.view.addSubview(textfield)
        
        let views: [String: AnyObject] = [
            "topLayoutGuide": self.topLayoutGuide,
            "textfield": textfield
        ]
        
        self.view.addConstraints(NSLayoutConstraint.constraintsWithVisualFormat("H:|-[textfield]-|", options: [], metrics: nil, views: views))
        self.view.addConstraints(NSLayoutConstraint.constraintsWithVisualFormat("V:|[topLayoutGuide]-50-[textfield]", options: [], metrics: nil, views: views))
    }

    override func didReceiveMemoryWarning() {
        super.didReceiveMemoryWarning()
        // Dispose of any resources that can be recreated.
    }
    
    func cancel() {
        self.dismissViewControllerAnimated(true, completion: nil)
    }
    
    func done() {
        guard let text = textfield.text where !text.isEmpty else {
            let alert = UIAlertController(title: "Invalid", message: "Please enter some text.", preferredStyle: .Alert)
            alert.addAction(UIAlertAction(title: "Cancel", style: .Cancel, handler: nil))
            self.presentViewController(alert, animated: true, completion: nil)
            return
        }
        let manager = CBLManager.sharedInstance()
        do {
            let database = try manager.databaseNamed("todos")
            let todo = GoItemNewTodo(text)
            todo.setCreatedAt(NSDate().timeIntervalSince1970 as Double)
            try Database.AddTodo(database, item: todo)
            self.dismissViewControllerAnimated(true, completion: nil)
        } catch {
            print(error)
        }
    }

}
