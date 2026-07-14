package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type AddUser struct {
	id   int
	name string
}

type AddConnection struct {
	fromID int
	toID   int
}

func TestAddUser(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    []AddUser
		expected []*User
	}{
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}}, expected: []*User{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}}},
		{input: []AddUser{{1, "Alice"}, {1, "Charlie"}}, expected: []*User{{ID: 1, Name: "Charlie"}}},
	}
	for _, tt := range tests {
		graph := NewGraph()
		for _, s := range tt.input {
			graph.AddUser(s.id, s.name)
		}
		assert.ElementsMatch(t, tt.expected, graph.GetAllUsers())

	}
}
func TestGetUser(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input        []AddUser
		expected     *User
		expectedBool bool
		userID       int
	}{
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}}, expected: &User{ID: 1, Name: "Alice"}, expectedBool: true, userID: 1},
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}}, expected: nil, expectedBool: false, userID: 0},
	}
	for _, tt := range tests {
		graph := NewGraph()
		for _, s := range tt.input {
			graph.AddUser(s.id, s.name)
		}
		user, result := graph.GetUser(tt.userID)
		assert.Equal(t, tt.expected, user)
		assert.Equal(t, tt.expectedBool, result)
	}
}

func TestAddConnection(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input           []AddUser
		expected        bool
		expectedCon     Connection
		expectedUsers   []*User
		id              []AddConnection
		fromID          int
		toID            int
		expectedReverse []*User
	}{
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}}, expected: true, expectedCon: Follower{}, expectedUsers: []*User{{ID: 2, Name: "Bob"}}, id: []AddConnection{{1, 2}}, fromID: 1, toID: 2, expectedReverse: []*User{}},
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}}, expected: false, expectedCon: nil, expectedUsers: []*User{}, id: []AddConnection{{1, 0}}, fromID: 1, toID: 0, expectedReverse: nil},
	}
	for _, tt := range tests {
		graph := NewGraph()
		for _, s := range tt.input {
			graph.AddUser(s.id, s.name)
		}
		var result bool
		for _, id := range tt.id {
			result = graph.AddConnection(id.fromID, id.toID)
		}
		getCon := graph.GetConnections(tt.fromID)
		reverse := graph.GetConnections(tt.toID)
		assert.ElementsMatch(t, tt.expectedReverse, reverse)
		assert.Equal(t, tt.expected, result)
		assert.Equal(t, tt.expectedCon, graph.connections[tt.fromID][tt.toID])
		assert.ElementsMatch(t, tt.expectedUsers, getCon)
	}
}

func TestGetConnections(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    []AddUser
		expected []*User
		userID   int
		id       []AddConnection
		fromID   int
		toID     int
	}{
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}}, expected: []*User{{ID: 2, Name: "Bob"}}, userID: 1, id: []AddConnection{{1, 2}}},
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}}, expected: nil, userID: 0, id: []AddConnection{{0, 2}}},
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}}, expected: []*User{}, userID: 2, id: []AddConnection{{1, 2}}},
	}
	for _, tt := range tests {
		graph := NewGraph()
		for _, s := range tt.input {
			graph.AddUser(s.id, s.name)
		}
		for _, id := range tt.id {
			graph.AddConnection(id.fromID, id.toID)
		}
		result := graph.GetConnections(tt.userID)
		assert.Equal(t, tt.expected, result)
	}
}

func TestHasConnection(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    []AddUser
		expected bool
		fromID   int
		toID     int
		id       []AddConnection
	}{
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}}, expected: true, fromID: 1, toID: 2, id: []AddConnection{{1, 2}}},
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}}, expected: false, fromID: 0, toID: 2, id: []AddConnection{{1, 2}}},
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}}, expected: false, fromID: 2, toID: 1, id: []AddConnection{{1, 2}}},
	}
	for _, tt := range tests {
		graph := NewGraph()
		for _, s := range tt.input {
			graph.AddUser(s.id, s.name)
		}
		for _, id := range tt.id {
			graph.AddConnection(id.fromID, id.toID)
		}
		result := graph.HasConnection(tt.fromID, tt.toID)
		assert.Equal(t, tt.expected, result)
	}
}

func TestUserCount(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    []AddUser
		expected int
	}{
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}}, expected: 2},
		{input: []AddUser{}, expected: 0},
	}
	for _, tt := range tests {
		graph := NewGraph()
		for _, s := range tt.input {
			graph.AddUser(s.id, s.name)
		}
		result := graph.UserCount()
		assert.Equal(t, tt.expected, result)
	}
}

func TestRemoveConnection(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input             []AddUser
		expected          bool
		fromID            int
		toID              int
		hasConAfterRemove bool
		id                []AddConnection
	}{
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}}, expected: true, fromID: 1, toID: 2, hasConAfterRemove: false, id: []AddConnection{{1, 2}}},
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}}, expected: false, fromID: 1, toID: 3, hasConAfterRemove: false, id: []AddConnection{{1, 2}}},
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}}, expected: false, fromID: 2, toID: 1, hasConAfterRemove: false, id: []AddConnection{{1, 2}}},
	}
	for _, tt := range tests {
		graph := NewGraph()
		for _, s := range tt.input {
			graph.AddUser(s.id, s.name)
		}
		for _, id := range tt.id {
			graph.AddConnection(id.fromID, id.toID)
		}
		result := graph.RemoveConnection(tt.fromID, tt.toID)
		assert.Equal(t, tt.expected, result)
		hasCon := graph.HasConnection(tt.fromID, tt.toID)
		assert.Equal(t, tt.hasConAfterRemove, hasCon)
	}
}

func TestRemoveUser(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input               []AddUser
		userID              int
		expected            bool
		hasConAfterRemove   bool
		expectedGetUserBool bool
		expectedGetUser     *User
		id                  []AddConnection
		fromID              int
		toID                int
	}{
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}}, userID: 1, expected: true, hasConAfterRemove: false, expectedGetUserBool: false, expectedGetUser: nil, id: []AddConnection{{1, 2}}, fromID: 1, toID: 2},
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}}, userID: 2, expected: true, hasConAfterRemove: false, expectedGetUserBool: false, expectedGetUser: nil, id: []AddConnection{{1, 2}}, fromID: 1, toID: 2},
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}}, userID: 0, expected: false, hasConAfterRemove: true, expectedGetUserBool: false, expectedGetUser: nil, id: []AddConnection{{1, 2}}, fromID: 1, toID: 2},
	}
	for _, tt := range tests {
		graph := NewGraph()
		for _, s := range tt.input {
			graph.AddUser(s.id, s.name)
		}
		for _, id := range tt.id {
			graph.AddConnection(id.fromID, id.toID)
		}
		result := graph.RemoveUser(tt.userID)
		assert.Equal(t, tt.expected, result)
		hasCon := graph.HasConnection(tt.fromID, tt.toID)
		assert.Equal(t, tt.hasConAfterRemove, hasCon)
		getUser, ok := graph.GetUser(tt.userID)
		assert.Equal(t, tt.expectedGetUserBool, ok)
		assert.Equal(t, tt.expectedGetUser, getUser)
	}
}

func TestIsMutual(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    []AddUser
		id       []AddConnection
		fromID   int
		toID     int
		expected bool
	}{
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}}, id: []AddConnection{{1, 2}, {2, 1}}, fromID: 1, toID: 2, expected: true},
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}, {3, "Charlie"}}, id: []AddConnection{{1, 2}, {1, 3}}, fromID: 1, toID: 2, expected: false},
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}}, id: []AddConnection{{3, 2}, {1, 2}}, fromID: 3, toID: 2, expected: false},
	}
	for _, tt := range tests {
		graph := NewGraph()
		for _, s := range tt.input {
			graph.AddUser(s.id, s.name)
		}
		for _, id := range tt.id {
			graph.AddConnection(id.fromID, id.toID)
		}
		result := graph.IsMutual(tt.fromID, tt.toID)
		assert.Equal(t, tt.expected, result)
	}
}

func TestConnectionCount(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    []AddUser
		id       []AddConnection
		userID   int
		expected int
	}{
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}, {3, "Charlie"}}, id: []AddConnection{{1, 2}, {1, 3}}, userID: 1, expected: 2},
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}, {3, "Charlie"}}, id: []AddConnection{{1, 2}, {1, 3}}, userID: 0, expected: 0},
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}, {3, "Charlie"}}, id: []AddConnection{}, userID: 1, expected: 0},
	}
	for _, tt := range tests {
		graph := NewGraph()
		for _, s := range tt.input {
			graph.AddUser(s.id, s.name)
		}
		for _, id := range tt.id {
			graph.AddConnection(id.fromID, id.toID)
		}
		result := graph.ConnectionCount(tt.userID)
		assert.Equal(t, tt.expected, result)
	}
}

func TestCommonConnections(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input       []AddUser
		id          []AddConnection
		userID1     int
		userID2     int
		expected    []*User
		afterBlock  []*User
		fromIDBlock int
		toIDBlock   int
	}{
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}, {3, "Charlie"}}, id: []AddConnection{{1, 3}, {3, 1}, {2, 3}, {3, 2}}, userID1: 1, userID2: 2, expected: []*User{{ID: 3, Name: "Charlie"}}, afterBlock: []*User{}, fromIDBlock: 1, toIDBlock: 3},
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}, {3, "Charlie"}}, id: []AddConnection{{1, 2}, {2, 1}}, userID1: 1, userID2: 2, expected: []*User{}, afterBlock: []*User{}},
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}, {3, "Charlie"}}, id: []AddConnection{{1, 3}, {3, 1}, {2, 3}, {3, 2}}, userID1: 0, userID2: 2, expected: nil, afterBlock: nil},
	}
	for _, tt := range tests {
		graph := NewGraph()
		for _, s := range tt.input {
			graph.AddUser(s.id, s.name)
		}
		for _, id := range tt.id {
			graph.AddConnection(id.fromID, id.toID)
		}
		result := graph.CommonConnections(tt.userID1, tt.userID2)
		assert.Equal(t, tt.expected, result)
		graph.AddTypedConnection(tt.fromIDBlock, tt.toIDBlock, Blocked{"spam"})
		result = graph.CommonConnections(tt.userID1, tt.userID2)
		assert.Equal(t, tt.afterBlock, result)
	}
}

func TestSuggestConnections(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input       []AddUser
		id          []AddConnection
		userID      int
		expected    []*User
		fromIDBlock int
		toIDBlock   int
		afterBlock  []*User
	}{
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}, {3, "Charlie"}}, id: []AddConnection{{1, 2}, {2, 3}}, userID: 1, expected: []*User{{ID: 3, Name: "Charlie"}}, afterBlock: []*User{{ID: 3, Name: "Charlie"}}},
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}, {3, "Charlie"}}, id: []AddConnection{{1, 2}, {2, 3}}, userID: 1, expected: []*User{{ID: 3, Name: "Charlie"}}, fromIDBlock: 1, toIDBlock: 2, afterBlock: []*User{}},
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}, {3, "Charlie"}}, id: []AddConnection{{1, 2}, {2, 3}}, userID: 1, expected: []*User{{ID: 3, Name: "Charlie"}}, fromIDBlock: 2, toIDBlock: 3, afterBlock: []*User{}},
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}, {3, "Charlie"}}, id: []AddConnection{{1, 2}, {2, 3}}, userID: 1, expected: []*User{{ID: 3, Name: "Charlie"}}, fromIDBlock: 1, toIDBlock: 3, afterBlock: []*User{}},
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}, {3, "Charlie"}}, id: []AddConnection{{1, 2}, {2, 3}}, userID: 0, expected: nil, afterBlock: nil},
	}
	for _, tt := range tests {
		graph := NewGraph()
		for _, s := range tt.input {
			graph.AddUser(s.id, s.name)
		}
		for _, id := range tt.id {
			graph.AddConnection(id.fromID, id.toID)
		}
		result := graph.SuggestConnections(tt.userID)
		assert.ElementsMatch(t, tt.expected, result)
		graph.AddTypedConnection(tt.fromIDBlock, tt.toIDBlock, Blocked{"spam"})
		result = graph.SuggestConnections(tt.userID)
		assert.ElementsMatch(t, tt.afterBlock, result)
	}
}

func TestGetAllUsers(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    []AddUser
		expected []*User
	}{
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}, {3, "Charlie"}}, expected: []*User{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}, {ID: 3, Name: "Charlie"}}},
		{input: nil, expected: nil},
		{input: []AddUser{}, expected: []*User{}},
	}

	for _, tt := range tests {
		graph := NewGraph()
		for _, s := range tt.input {
			graph.AddUser(s.id, s.name)
		}
		result := graph.GetAllUsers()
		assert.ElementsMatch(t, tt.expected, result)
	}
}

func TestAddTypedConnection(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    []AddUser
		id       []AddConnection
		conn     Connection
		fromID   int
		toID     int
		expected bool
	}{
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}}, id: []AddConnection{{1, 2}}, conn: Follower{true}, fromID: 1, toID: 2, expected: true},
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}}, id: []AddConnection{}, conn: Friend{"01012026"}, fromID: 1, toID: 2, expected: true},
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}}, id: []AddConnection{{1, 2}}, conn: Blocked{"spam"}, fromID: 2, toID: 1, expected: true},
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}}, id: []AddConnection{}, conn: nil, fromID: 0, toID: 2, expected: false},
	}
	for _, tt := range tests {
		graph := NewGraph()
		for _, s := range tt.input {
			graph.AddUser(s.id, s.name)
		}
		for _, id := range tt.id {
			graph.AddConnection(id.fromID, id.toID)
		}
		result := graph.AddTypedConnection(tt.fromID, tt.toID, tt.conn)
		assert.Equal(t, tt.expected, result)
		connInfo, _ := graph.GetConnectionInfo(tt.fromID, tt.toID)
		assert.Equal(t, tt.conn, connInfo)
	}
}

func TestGetConnectionsByType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    []AddUser
		id       []AddConnection
		connType string
		conn     Connection
		fromID   int
		toID     int
		userID   int
		expected []*User
	}{
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}}, id: []AddConnection{{1, 2}}, connType: "following", conn: Follower{true}, fromID: 1, toID: 2, userID: 1, expected: nil},
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}}, id: []AddConnection{{1, 2}}, connType: "friend", conn: Friend{"01012026"}, fromID: 1, toID: 2, userID: 1, expected: []*User{{ID: 2, Name: "Bob"}}},
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}}, id: []AddConnection{{1, 2}}, connType: "friend", conn: Blocked{"spam"}, fromID: 1, toID: 2, userID: 0, expected: nil},
	}
	for _, tt := range tests {
		graph := NewGraph()
		for _, s := range tt.input {
			graph.AddUser(s.id, s.name)
		}
		for _, id := range tt.id {
			graph.AddConnection(id.fromID, id.toID)
		}
		graph.AddTypedConnection(tt.fromID, tt.toID, tt.conn)
		result := graph.GetConnectionsByType(tt.userID, tt.connType)
		assert.ElementsMatch(t, tt.expected, result)
	}
}

func TestGetConnectionInfo(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input        []AddUser
		id           []AddConnection
		conn         Connection
		fromID       int
		toID         int
		expected     Connection
		expectedBool bool
	}{
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}}, id: []AddConnection{{1, 2}}, conn: Follower{true}, fromID: 1, toID: 2, expected: Follower{true}, expectedBool: true},
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}}, id: []AddConnection{}, conn: Friend{"01012026"}, fromID: 1, toID: 2, expected: Friend{"01012026"}, expectedBool: true},
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}}, id: []AddConnection{{1, 2}}, conn: Blocked{"spam"}, fromID: 2, toID: 1, expected: Blocked{"spam"}, expectedBool: true},
		{input: []AddUser{{1, "Alice"}, {2, "Bob"}}, id: []AddConnection{}, conn: Blocked{"spam"}, fromID: 0, toID: 2, expected: nil, expectedBool: false},
	}
	for _, tt := range tests {
		graph := NewGraph()
		for _, s := range tt.input {
			graph.AddUser(s.id, s.name)
		}
		for _, id := range tt.id {
			graph.AddConnection(id.fromID, id.toID)
		}
		graph.AddTypedConnection(tt.fromID, tt.toID, tt.conn)
		connInfo, ok := graph.GetConnectionInfo(tt.fromID, tt.toID)
		assert.Equal(t, tt.expected, connInfo)
		assert.Equal(t, tt.expectedBool, ok)
	}
}
