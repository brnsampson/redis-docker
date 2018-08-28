package instance

import "testing"

// Implement mock for third party library shim while testing.
type TestInstanceShim struct{
	err error
}

func(s *TestInstanceShim) Ping() error {
	return c.err
}

func(s *TestInstanceShim) GetInfo() (string, error) {
	return "info", c.err
}

func(s *TestInstanceShim) SetReplication(master string) error {
	return c.err
}

func (c *TestInstanceShim) ReadConfig(k string) (string, error) {
	return k, c.err
}

func (c *TestInstanceShim) UpdateConfig(k, v string) error {
	return c.err
}

func(s *TestInstanceShim) Quit() error {
	return c.err
}

// Mock out InstanceInfoParser for Instance testing

type TestInfoParser struct{
	err error
	role string
	loading string
	master_link_status string
	master_sync_left_bytes string
	master_last_io_seconds string
	master_host string
	master_port string
}

func(ip *TestInfoParser) ParseInfo(_ string, ii *InstanceInfo) error {
	if ip.role != "" {
		ii["role"] = ip.role
	}

	if ip.loading != "" {
		ii["loading"] = ip.loading
	}

	if ip.master_link_status != "" {
		ii["master_link_status"] = ip.master_link_status
	}

	if ip.master_last_io_seconds != "" {
		ii["master_last_io_seconds"] = ip.master_last_io_seconds
	}

	if ip.master_sync_left_bytes != "" {
		ii["master_sync_left_bytes"] = ip.master_sync_left_bytes
	}

	if ip.master_host != "" {
		ii["master_host"] = ip.master_host
	}

	if ip.master_port != "" {
		ii["master_port"] = ip.master_port
	}
	return ip.err
}

// Perform tests
func TestParseInfo(t *testing.T) {
	ip := InfoParser{}

	// Test InfoParser for a valid master info string.
	rawInfo := `# Server
redis_version:4.0.11
port:6379

# Persistence
loading:0

# Replication
role:master
connected_slaves:0
`
	ii := make(InstanceInfo)
	err := ip.ParseRedisInfo(rawInfo, &ii)
	if err {
		t.Errorf("Recieved unexpected error from master info %q", err)
	}
	if ii["role"] != "master" || ii["loading"] != "0" {
		t.Error("did not recieve expected result. Instead got %q", ii)
	}

	// Test InfoParser for a valid replica info string.
	rawInfo = `# Server
redis_version:4.0.11
port:6379

# Persistence
loading:0

# Replication
role:slave
connected_slaves:0
slave_repl_offset:0
master_sync_in_progress:0
master_link_status:up
master_last_io_seconds_ago:0
master_sync_in_progress:0
`
	err = ip.ParseRedisInfo(rawInfo, &ii)
	if err {
		t.Errorf("Recieved unexpected error from replica info %q", err)
	}
	if ii["role"] != "slave" || ii["loading"] != "0" || ii["master_sync_in_progress"] != "0" || ii["master_link_status"] != "up" {
		t.Error("did not recieve expected result. Instead got %q", ii)
	}

	// Test InfoParser for a replicating info string.
	rawInfo = `# Server
redis_version:4.0.11
port:6379

# Persistence
loading:0

# Replication
role:slave
connected_slaves:0
slave_repl_offset:0
master_link_status:up
master_last_io_seconds_ago:0
master_sync_in_progress:1
`
	err = ip.ParseRedisInfo(rawInfo, &ri)
	if err {
		t.Errorf("Recieved unexpected error from replicating info %q", err)
	}
	if ii["role"] != "slave" || ii["loading"] != "0" || ii["master_sync_in_progress"] != "1" || ii["master_link_status"] != "up" {
		t.Error("did not recieve expected result. Instead got %q", ii)
	}
}

func TestIsInstanceReady(t *testing.T) {
	// Test master case
	ts := TestInstanceShim{err: nil}
	ip := TestInstanceParser{
		err: nil,
		role: "master",
		loading: 0,
	}

	inst := Instance{&ts, &ip}

	ready := inst.IsRedisReady()
	if !ready {
		t.Error("Loaded master redis should return ready")
	}

	// Test replica with no lag
	ip = InfoParser{
		err: nil,
		role: "slave",
		loading: "0",
		master_link_status "up",
		master_sync_left_bytes: "0",
	}
	inst = Instance{&ts, &ip}

	ready = inst.IsRedisReady()
	if !ready {
		t.Error("Loaded replica redis should return ready")
	}

	// Test instance which is still loading
	ip = InfoParser{
		err: nil,
		role: "master",
		loading: "1",
	}
	inst = Instance{&ts, &ip}

	ready = inst.IsRedisReady()
	if ready {
		t.Error("Loading redis instance should not return ready")
	}

	// Test replica still loading from master
	ip = InfoParser{
		err: nil,
		role: "slave",
		loading: "0",
		master_link_status "up",
		master_sync_left_bytes: "402851",
		master_last_to_seconds: "0",
	}
	inst = Instance{&ts, &ip}

	ready = inst.IsRedisReady()
	if ready {
		t.Error("Replicating redis replica should not return ready")
	}

	ip = InfoParser{
		err: fmt.Error("test error for IsInstanceReady"),
	}
	inst = Instance{&ts, &ip}
	ready, err = inst.IsRedisReady()
	if err != fmt.Error("test error for IsInstanceReady") {
		t.Error("Did not recieve error from ParseInfo in IsInstanceReady")
	}
}

func TestReadInstanceMaster(t *testing.T) {
	// Test a master instance
	ts := TestInstanceShim{err: nil}
	ip := TestInstanceParser{
		err: nil,
		role: "master",
		loading: 0,
	}

	inst := Instance{&ts, &ip}
	master, err := inst.ReadInstanceMaster()

	if master != "127.0.0.1:port" {
		t.Error("master instances should return a port on localhost"
	}

	// Test a slave instance
	ip = InfoParser{
		err: nil,
		role: "slave",
		master_host: "192.168.0.1",
		master_port: "12345",
		loading: "0",
		master_link_status "up",
		master_sync_left_bytes: "402851",
		master_last_to_seconds: "0",
	}
	inst = Instance{&ts, &ip}
	master, err = inst.ReadInstanceMaster()

	if master != "192.168.0.1:12345" {
		t.Error("master instance should be constructed from parsed info for slave"
	}

	ip = TestInstanceParser{
		err: fmt.Error("error in ParseInstanceInfo"),
		role: "master",
		loading: 0,
	}
	inst = Instance{&ts, &ip}
	master, err = inst.ReadInstanceMaster()

	if err != fmt.Error("error in ParseInstanceInfo") {
		t.Error("Error returned by ParseInstanceInfo should propogate up")
	}
}

func TestUpdateInstanceMaster(t *testing.T) {
	ts := TestInstanceShim{err: nil}
	ip := TestInstanceParser{
		err: nil,
		role: "master",
		loading: 0,
	}

	inst := Instance{&ts, &ip}
	err := inst.UpdateInstanceMaster("192.168.0.1:12345")
	if err == nil {
		t.Error("Unsuccessful update from UpdateInstanceMaster does not produce error")
	}

	// Test error in the InstanceShimmer
	ts = TestInstanceShim{err: fmt.Error("error in TestInstanceShim")}
	inst = Instance{&ts, &ip}
	err = inst.UpdateInstanceMaster()
	if err != fmt.Error("error in TestInstanceShim") {
		t.Error("Error from TestInstanceShimmer not raised through UpdateInstanceMaster")
	}

	// Test error in the parser
	ts = TestInstanceShim{err: nil}
	ip = TestInstanceParser{
		err: fmt.Error("error in parser"),
		role: "master",
		loading: 0,
	}
	inst = Instance{&ts, &ip}

	err = inst.UpdateInstanceMaster("192.168.0.1:12345")
	if err != fmt.Error("error in parser") {
		t.Error("Error from TestInstanceParser not raised through UpdateInstanceMaster")
	}

	// Test what should work
	ip = InfoParser{
		err: nil,
		role: "slave",
		master_host: "192.168.0.1",
		master_port: "12345",
		loading: "0",
		master_link_status "up",
		master_sync_left_bytes: "402851",
		master_last_to_seconds: "0",
	}
	inst = Instance{&ts, &ip}

	err = inst.UpdateInstanceMaster("192.168.0.1:12345")
	if err != nil
		t.Error("Got unexpected error when succesfully updating master")
	}
}

func TestClaimInstanceMaster(t *testing.T) {
	// Test successful case
	ts := TestInstanceShim{err: nil}
	ip := TestInstanceParser{
		err: nil,
		role: "master",
		loading: 0,
	}

	inst := Instance{&ts, &ip}
	err := inst.ClaimInstanceMaster()
	if err != nil {
		t.Error("Error recieved for successful case")
	}

	// Test error in the InstanceShimmer
	ts = TestInstanceShim{err: fmt.Error("error in TestInstanceShim")}
	inst = Instance{&ts, &ip}
	err = inst.UpdateInstanceMaster()
	if err != fmt.Error("error in TestInstanceShim") {
		t.Error("Error from TestInstanceShimmer not raised through ClaimInstanceMaster")
	}

	// Test error in the parser
	ts = TestInstanceShim{err: nil}
	ip = TestInstanceParser{
		err: fmt.Error("error in parser"),
		role: "master",
		loading: 0,
	}
	inst = Instance{&ts, &ip}

	err = inst.ClaimInstanceMaster()
	if err != fmt.Error("error in parser") {
		t.Error("Error from TestInstanceParser not raised through ClaimInstanceMaster")
	}

	// Test failure to update master
	ip = InfoParser{
		err: nil,
		role: "slave",
		master_host: "192.168.0.1",
		master_port: "12345",
		loading: "0",
		master_link_status "up",
		master_sync_left_bytes: "402851",
		master_last_to_seconds: "0",
	}
	inst = Instance{&ts, &ip}

	err = inst.ClaimInstanceMaster()
	if err == nil
		t.Error("Failed to return error when update to master failed")
	}
}
