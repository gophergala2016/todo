
import UIKit

class AddTodoViewController: UIViewController {
    
    private var textfield: UITextField!

    override func viewDidLoad() {
        super.viewDidLoad()

        self.title = "Add todo"
        self.view.backgroundColor = UIColor.whiteColor()
        self.navigationItem.leftBarButtonItem = UIBarButtonItem(barButtonSystemItem: .Cancel, target: self, action: "cancel")
        
        // add textfield
        textfield = UITextField()
        textfield.translatesAutoresizingMaskIntoConstraints = false
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

}
