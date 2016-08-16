import optparse, os, queue, sys, threading, tkinter

gfx_dict = dict()


def parse_args():

	global opts

	opt_parser = optparse.OptionParser()

	opt_parser.add_option(
		"--width",
		dest = "width",
		type = "int",
		help = "Width in characters [default: %default]")
	opt_parser.set_defaults(width = 800)

	opt_parser.add_option(
		"--height",
		dest = "height",
		type = "int",
		help = "Height in characters [default: %default]")
	opt_parser.set_defaults(height = 600)

	opt_parser.add_option(
		"--directory",
		dest = "directory",
		type = "str",
		help = "GFX directory [default: %default]")
	opt_parser.set_defaults(charfile = ".")

	opt_parser.add_option(
		"--bg",
		dest = "bg",
		type = "str",
		help = "Background colour [default: %default]")
	opt_parser.set_defaults(bg = "white")

	opts, __ = opt_parser.parse_args()


def load_graphics(directory):
	for dirpath, dirnames, filenames in os.walk(directory):
		for filename in filenames:
			fullpath = os.path.join(dirpath, filename)
			if fullpath.lower().endswith(".gif"):
				gfx_dict[filename] = tkinter.PhotoImage(file = fullpath)


def input_thread(q):
	while 1:
		s = input()
		q.put(s)


class Renderer(tkinter.Canvas):
	def __init__(self, owner, *args, **kwargs):
		tkinter.Canvas.__init__(self, owner, *args, **kwargs)

		self.input_queue = queue.Queue()
		threading.Thread(target = input_thread, daemon = True, kwargs = {"q" : self.input_queue}).start()

		self.poller()

	def poller(self):

		frame_started = False

		# Once we've started receiving a frame, we need to continue until it ends

		while 1:
			try:
				s = self.input_queue.get(block = False)

				if not frame_started:
					self.delete(tkinter.ALL)
					frame_started = True

				fields = s.split()

				if len(fields) == 3:
					spritename, x, y = fields[0], int(fields[1]), int(fields[2])
					self.create_image(x, y, image = gfx_dict[spritename])

				if len(fields) == 1:
					if fields[0] == "ENDFRAME":
						self.update_idletasks()
						frame_started = False
						print("ENDFRAME")
						sys.stdout.flush()

			except queue.Empty:
				if not frame_started:
					break

		self.after(1, self.poller)


class Root(tkinter.Tk):
	def __init__(self, *args, **kwargs):

		parse_args()

		tkinter.Tk.__init__(self, *args, **kwargs)
		self.resizable(width = False, height = False)

		load_graphics(opts.directory)

		virtue = Renderer(self, width = opts.width, height = opts.height, bd = 0, highlightthickness = 0, bg = opts.bg)
		virtue.pack()



if __name__ == "__main__":
	app = Root()
	app.mainloop()
