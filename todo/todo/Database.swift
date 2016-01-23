
import Foundation
import Item

public class Database {

    public class func GetTodos(db: CBLDatabase) throws -> [GoItemTodo] {
        let view = db.viewNamed("todos")
        // use basically the same map function as in javascript
        let map: CBLMapBlock = { doc, emit in
            if
                let type = doc["type"] as? String where type == "todo",
                let createdAt = doc["createdAt"] as? Float64 {
                    emit(createdAt, nil)
            }
        }
        view.setMapBlock(map, version: "2")
        let query = db.viewNamed("todos").createQuery()
        query.startKey = 1356998400 as Float64
        query.endKey = 32503593600 as Float64
        do {
            let data = try query.run()
            var todos: [GoItemTodo] = []
            for var index = 0; index < Int(data.count); ++index {
                if let doc = data.rowAtIndex(UInt(index)).document {
                    if
                        let type = doc["type"] as? String,
                        let text = doc["text"] as? String,
                        let createdAt = doc["createdAt"] as? Float64,
                        let done = doc["done"] as? Bool {
                            let item = GoItemNewTodo(text)
                            item.setType(type)
                            item.setCreatedAt(createdAt)
                            item.setDone(done)
                            todos.append(item)
                    }
                    
                }
            }
            return todos
        } catch {
            throw error
        }
    }
    
}