tw:
	tailwindcss -i ./static/css/input.css -o ./static/css/output.css --watch

install:
	flint spark; mv out ..; cd ..; mv out www.phillip-england.com; cd www.phillip-england.com; rm -r index.html; rm -r post; rm -r posts.html; rm -r static; cd out; mv * ..; cd ..; rm -r out;