
const mapper = function(){ 
    return {
        $el: null,
        $canvas: null,
        room_click: null,
        new: function($el, clickHandler) {
            this.$el = $el;
            let $canvas = document.createElement("canvas");
            $canvas.clientWidth = $el.clientWidth;
            $canvas.clientHeight = $el.clientHeight;
            this.$canvas = $canvas;
            this.$el.appendChild(this.$canvas);
            this.room_click = clickHandler;
            return this;
        },
        show: function(area, room) {
            this.$el.innerHTML="";
            this.$el.appendChild(this.$canvas);
            let visited = [];
            const ctx = this.$canvas.getContext('2d');
            const w = this.$el.clientWidth;
            const h = this.$el.clientHeight;
            let startX = w/2;
            let startY = h/2;
            let lastX = 0;
            let lastY = 0;
            ctx.canvas.width = w;
            ctx.canvas.height = h;
            ctx.fillStyle = "grey";
            ctx.strokeStyle = "white";
            ctx.clearRect(0, 0, w, h);
            ctx.moveTo(startX, startY);
            const gap = 45;
            let $this = this;
            let get_room = function(roomId) {
                for(const room of area.rooms) {
                    if (room.id == roomId) {
                        return room;
                    }
                }
                return null;
            }
            let draw_exit = function(ctx, x, y, dx, dy) {
                ctx.beginPath();
                ctx.moveTo(x, y);
                ctx.lineTo(dx, dy);
                ctx.closePath();
                ctx.stroke();
            }
            let draw_room = function(ctx, room, x, y) {
                if (!visited.includes(room.id)){
                    visited.push(room.id);
                    ctx.fillStyle = "white";
                    ctx.fillRect(x - (gap/2), y - (gap/2), gap, gap);
                    let $room = document.createElement("div");
                    $room.style.position = "absolute";
                    $room.style.left = (x - (gap/2)).toString()+"px";
                    $room.style.top = (y - (gap/2)).toString()+"px";
                    $room.style.width = gap.toString()+"px";
                    $room.style.height = gap.toString()+"px";
                    $room.clientWidth = gap.toString()+"px";
                    $room.clientHeight = gap.toString()+"px";
                    $room.style.background = "#ffffff";
                    $room.style.color = "#FF9933";
                    $room.innerHTML = room.id;
                    $room.onclick = function(event){ event.preventDefault(); $this.room_click(room); };
                    $this.$el.appendChild($room);
                    if (room.exits) {
                        if (Object.keys(room.exits).length > 0) {
                            draw_exits(ctx, room, x, y);
                        }
                    }
                }
            }
            let draw_exits = function(ctx, room, x, y) {
                for(const [exit, roomId] of Object.entries(room.exits)){
                    switch (exit) {
                        case "north":
                            draw_exit(ctx, x, y, x, y-gap);
                            draw_room(ctx, get_room(roomId), x, y-gap-gap);
                            break;
                        case "south":
                            draw_exit(ctx, x, y, x, y+gap);
                            draw_room(ctx, get_room(roomId), x, y+gap+gap);
                            break;
                        case "east":
                            draw_exit(ctx, x, y, x+gap, y);
                            draw_room(ctx, get_room(roomId), x+gap+gap, y);
                            break;
                        case "west":
                            draw_exit(ctx, x, y, x-gap, y);
                            draw_room(ctx, get_room(roomId), x-gap-gap, y);
                            break;
                        case "northeast":
                            draw_exit(ctx, x, y, x+gap, y-gap);
                            draw_room(ctx, get_room(roomId), x+gap+gap, y-gap-gap);
                            break;
                        case "southeast":
                            draw_exit(ctx, x, y, x+gap, y+gap);
                            draw_room(ctx, get_room(roomId), x+gap+gap, y+gap+gap);
                            break;
                        case "southwest":
                            draw_exit(ctx, x, y, x-gap, y+gap);
                            draw_room(ctx, get_room(roomId), x-gap-gap, y+gap+gap);
                            break;
                        case "northwest":
                            draw_exit(ctx, x, y, x-gap, y-gap);
                            draw_room(ctx, get_room(roomId), x-gap-gap, y-gap-gap);
                            break;
                    }
                }
                lastX = x;
                lastY = y;
            }
            if (room != null){
                draw_room(ctx, room, startX, startY);
            }
            
        }
    }
};