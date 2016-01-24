
import Foundation
import Item

public class Database {
    
    // get database for main thread
    public class func Get() throws -> CBLDatabase {
        let manager = CBLManager.sharedInstance()
        do {
            return try manager.databaseNamed("todos")
        } catch {
            throw error
        }
    }
    
    public class func GetTodoByID(db: CBLDatabase, id: String) -> GoItemTodo {
        let doc = db.documentWithID(id)
        if let doc = doc {
            if
                let id = doc["_id"] as? String,
                let rev = doc["_rev"] as? String,
                let type = doc["type"] as? String,
                let text = doc["text"] as? String,
                let createdAt = doc["createdAt"] as? Float64,
                let done = doc["done"] as? Bool {
                    let item = GoItemNewTodo(text)
                    item.setID(id)
                    item.setRev(rev)
                    item.setType(type)
                    item.setCreatedAt(createdAt)
                    item.setDone(done)
                    return item
            }
        }
        return GoItemNewTodo("")
    }
    
    public class func AddTodo(db: CBLDatabase, item: GoItemTodo) throws {
        do {
            let document = db.createDocument()
            let properties: [String: AnyObject] = [
                "type": item.type(),
                "text": item.text(),
                "createdAt": item.createdAt(),
                "done": item.done()
            ]
            try document.putProperties(properties)
        } catch {
            throw error
        }
    }
    
    public class func UpdateTodo(db: CBLDatabase, item: GoItemTodo) throws {
        let document = db.documentWithID(item.id())
        do {
            try document!.update({ rev -> Bool in
                rev["type"] = item.type()
                rev["text"] = item.text()
                rev["createdAt"] = item.createdAt()
                rev["done"] = item.done()
                return true
            })
        } catch {
            throw error
        }
    }
    
    public class func RemoveTodoByID(db: CBLDatabase, id: String) throws {
        let document = db.documentWithID(id)
        do {
            try document!.deleteDocument()
        } catch {
            throw error
        }
    }

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
        view.setMapBlock(map, version: "1")
        let query = db.viewNamed("todos").createQuery()
        query.startKey = 1356998400 as Float64
        query.endKey = 32503593600 as Float64
        do {
            let data = try query.run()
            var todos: [GoItemTodo] = []
            for var index = 0; index < Int(data.count); ++index {
                if let doc = data.rowAtIndex(UInt(index)).document {
                    if
                        let id = doc["_id"] as? String,
                        let rev = doc["_rev"] as? String,
                        let type = doc["type"] as? String,
                        let text = doc["text"] as? String,
                        let createdAt = doc["createdAt"] as? Float64,
                        let done = doc["done"] as? Bool {
                            let item = GoItemNewTodo(text)
                            item.setID(id)
                            item.setRev(rev)
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