package tiles

import "os"
import "encoding/binary"
import "bufio"
import "http"
import "strconv"

import "fmt"

type TileServer struct {
	size int32; // Number of tiles.
	x, y []int32; // coordinates of tile images
	zoom []int8;  // corresponding zooms
	start []int32; // corresponding offsets into data file
	data *os.File;
}

func Load(name string) (ts *TileServer, err os.Error) {
	set := new(TileServer);
	index_file, err := os.Open(name + "/index.new", os.O_RDONLY, 0);
	if err != nil {
		return nil, err;
	}

	set.data, err = os.Open(name + "/map.data", os.O_RDONLY, 0);
	if err != nil {
		return nil, err;
	}

	index := bufio.NewReader(index_file);
	binary.Read(index, binary.BigEndian, &set.size);

	// Debug
	fmt.Println("size=", set.size);
	dir, _ := os.Stat(name + "/index.new");
	fmt.Println("filesize=", dir.Size);
	fmt.Println("expected=", set.size * 4 * 3 + set.size + 4 + 4);
	set.x = make([]int32, set.size);
	set.y = make([]int32, set.size);
	set.zoom = make([]int8, set.size);
	set.start = make([]int32, set.size+1);


	err = binary.Read(index, binary.BigEndian, set.x);
	if err == nil { err = binary.Read(index, binary.BigEndian, set.y); }
	if err == nil { err = binary.Read(index, binary.BigEndian, set.zoom); }
	if err == nil { err = binary.Read(index, binary.BigEndian, set.start); }

	if err != nil {
		return nil, err;
	}

	return set, nil;
}


func (set *TileServer) compare(zoom, x, y, ind int) int {
	if zoom != int(set.zoom[ind]) {
		return zoom - int(set.zoom[ind]);
	}
	if x != int(set.x[ind]) {
		return x - int(set.x[ind]);
	}
	return y - int(set.y[ind]);
}

func (set *TileServer) search(zoom, x, y int) int {
	low := 0;
	high := int(set.size - 1);
	for  high - low > 1 {
		middle := (low + high) / 2;
		d := set.compare(zoom, x, y, middle);
		if d == 0 {
			return middle;
		}
		if d < 0 {
			high = middle;
		} else {
			low = middle + 1;
		}
	}
	if set.compare(zoom, x, y, low) == 0 {
		return low;
	}
	return -1;
}

func (set *TileServer) GetTile(x, y, zoom int) []byte {
	index := set.search(zoom, x, y);
	if index < 0 {
		return nil;
	}
	offset := set.start[index];
	length := set.start[index+1] - offset;
	buffer := make([]byte, length);
	n, err := set.data.ReadAt(buffer, int64(offset));
	if n != int(length) {
		fmt.Println("Incomplete read");
		return nil;
	}
	if err != nil {
		fmt.Println("Err:", err.String());
		return nil;
	}
		
	return buffer;
}

func (set *TileServer) ServeHTTP(conn *http.Conn, req *http.Request) {
	req.ParseForm();
	var x, y, zoom int;
	x, err := strconv.Atoi(req.Form["x"][0]);
	if err == nil { y, err = strconv.Atoi(req.Form["y"][0]); }
	if err == nil { zoom, err = strconv.Atoi(req.Form["z"][0]); }

	if err != nil {
		fmt.Println("Err:", err.String());
		conn.WriteHeader(404);
		return;
	}

	tile := set.GetTile(x, y, zoom);
	if tile == nil {
		fmt.Println("No tile at", x, y, zoom);
		conn.WriteHeader(404);
		return;
	}
	conn.Write(tile);
}
