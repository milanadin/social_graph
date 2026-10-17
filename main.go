package main

import (
	"fmt"
	"maps"
	"slices"
)

type User struct {
	ID   int
	Name string
}

type Graph struct {
	users       map[int]*User
	connections map[int]map[int]Connection
}

func NewGraph() *Graph {
	users := make(map[int]*User)
	connections := make(map[int]map[int]Connection)
	newGraph := &Graph{users: users, connections: connections}
	return newGraph
}

func (g *Graph) AddUser(id int, name string) {
	if _, ok := g.users[id]; ok {
		g.users[id].Name = name
		return
	}
	g.users[id] = &User{ID: id, Name: name}

}

func (g *Graph) GetUser(id int) (*User, bool) {
	users, ok := g.users[id]
	if !ok {
		return nil, ok
	}
	return users, ok
}

func (g *Graph) AddConnection(fromID, toID int) bool {
	_, okFrom := g.users[fromID]
	_, okTo := g.users[toID]
	if !okFrom || !okTo {
		return false
	}
	_, ok := g.connections[fromID]
	if !ok {
		g.connections[fromID] = map[int]Connection{}
	}
	g.connections[fromID][toID] = Follower{}
	return true
}

func (g *Graph) GetConnections(userID int) []*User {
	_, ok := g.users[userID]
	if !ok {
		return nil
	}
	m := g.connections[userID]
	sliseUsers := make([]*User, 0, len(g.connections[userID]))
	for id := range m {
		sliseUsers = append(sliseUsers, g.users[id])
	}
	return sliseUsers
}

func (g *Graph) HasConnection(fromID, toID int) bool {
	// TODO
	_, okFrom := g.users[fromID]
	_, okTo := g.users[toID]
	if !okFrom || !okTo {
		return false
	}
	_, ok := g.connections[fromID][toID]
	return ok
}

func (g *Graph) UserCount() int {
	return len(g.users)
}

func (g *Graph) RemoveConnection(fromID, toID int) bool {
	if g.HasConnection(fromID, toID) {
		delete(g.connections[fromID], toID)
		return true
	}
	return false
}

func (g *Graph) RemoveUser(id int) bool {
	if _, ok := g.users[id]; !ok {
		return false
	}

	for k := range g.connections {
		if g.HasConnection(k, id) {
			delete(g.connections[k], id)
		}
	}
	delete(g.users, id)
	delete(g.connections, id)
	return true
}

func (g *Graph) IsMutual(id1, id2 int) bool {
	return g.HasConnection(id1, id2) && g.HasConnection(id2, id1)
}

func (g *Graph) ConnectionCount(userID int) int {
	return len(g.connections[userID])
}

func (g *Graph) CommonConnections(id1, id2 int) []*User {
	_, ok1 := g.users[id1]
	_, ok2 := g.users[id2]
	if !ok1 || !ok2 {
		return nil
	}
	sliceUsers := make([]*User, 0, len(g.connections[id1]))
	sliceUsersID1 := g.GetConnections(id1)
	for _, s := range sliceUsersID1 {
		getConn1, con1 := g.GetConnectionInfo(id1, s.ID)
		getConn2, con2 := g.GetConnectionInfo(id2, s.ID)
		if g.HasConnection(s.ID, id1) && g.IsMutual(id2, s.ID) && con1 && con2 && getConn1.Type() != "blocked" && getConn2.Type() != "blocked" {
			sliceUsers = append(sliceUsers, s)
		}

	}
	return sliceUsers
}

func (g *Graph) SuggestConnections(userID int) []*User { // возвращает связи, кроме заблокированных пользователей и тех, кого заблокировали подписки пользователя
	if _, ok := g.users[userID]; !ok {
		return nil
	}

	mapConnections := g.connections[userID] // все связи пользователя
	mapConCon := make(map[int]Connection)   // все связи друзей и подписок
	for k, v := range mapConnections {
		if v.Type() == "blocked" {
			continue // пропуск заблокированных
		}
		maps.Copy(mapConCon, g.connections[k])
	}

	for k, v := range mapConCon {
		con, ok := g.connections[userID][k]
		if v.Type() == "blocked" {
			delete(mapConCon, k)
		}
		if ok && con.Type() == "blocked" { // проверка предложения подписки, которую пользователь заблокировал
			delete(mapConCon, k)

		}
	}
	delete(mapConCon, userID)
	for alreadySub := range mapConnections {
		delete(mapConCon, alreadySub)
	}
	suggestConnections := make([]*User, 0, len(mapConCon))
	for k := range mapConCon {
		if user, ok := g.users[k]; ok {
			suggestConnections = append(suggestConnections, user)
		}
	}
	return suggestConnections
}

func (g *Graph) GetAllUsers() []*User {
	return slices.Collect(maps.Values(g.users))
}

type Connection interface {
	Type() string
	Weight() int
}

type Friend struct {
	Since string // Дата начала дружбы
}

func (f Friend) Type() string {
	return "friend"
}

func (f Friend) Weight() int {
	return 10
}

type Follower struct {
	Notifications bool
}

func (f Follower) Type() string {
	return "follower"
}

func (f Follower) Weight() int {
	return 5
}

type Blocked struct {
	Reason string
}

func (b Blocked) Type() string {
	return "blocked"
}

func (b Blocked) Weight() int {
	return -1
}

func (g *Graph) AddTypedConnection(fromID, toID int, conn Connection) bool {
	_, okFrom := g.users[fromID]
	_, okTo := g.users[toID]
	if !okFrom || !okTo {
		return false
	}
	_, ok := g.connections[fromID]
	if !ok {
		g.connections[fromID] = map[int]Connection{}
	}

	g.connections[fromID][toID] = conn
	return true
}

func (g *Graph) GetConnectionsByType(userID int, connType string) []*User {
	_, ok := g.users[userID]
	if !ok || (connType != "friend" && connType != "blocked" && connType != "follower") {
		return nil
	}
	mUsers := g.connections[userID]
	connectionsByType := make([]*User, 0, len(mUsers))
	for k, v := range mUsers {
		if v.Type() == connType {
			connectionsByType = append(connectionsByType, g.users[k])
		}
	}
	return connectionsByType
}

func (g *Graph) GetConnectionInfo(fromID, toID int) (Connection, bool) {
	_, okFrom := g.users[fromID]
	_, okTo := g.users[toID]
	if !okFrom || !okTo {
		return nil, false
	}
	connection := g.connections[fromID][toID]
	if connection == nil {
		return nil, false
	}
	return connection, true
}

func main() {
	graph := NewGraph()

	graph.AddUser(1, "Alice")
	graph.AddUser(2, "Bob")
	graph.AddUser(3, "Charlie")

	graph.AddConnection(1, 2) // Alice -> Bob
	graph.AddConnection(3, 1) // Charlie -> Alice
	graph.AddConnection(3, 2) // Charlie -> Bob
	graph.AddConnection(2, 3) // Bob -> Charlie
	graph.AddConnection(2, 1) // Bob -> Charlie

	if user, ok := graph.GetUser(1); ok {
		fmt.Printf("User: %s\n", user.Name)
		friends := graph.GetConnections(1)
		fmt.Printf("Friends: %d\n", len(friends))
		for _, friend := range friends {
			fmt.Printf("  - %s\n", friend.Name)
		}
	}

	fmt.Printf("Alice and Bob connected: %v\n",
		graph.HasConnection(1, 2))

}
