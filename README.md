# BAE - brainsik attractor explorer

This is the project I'm using to learn Go. As such, it's focus is trying and learning the language more than being usable by anyone else.

Here's [some images](https://hachyderm.io/@brainsik/110766510259586543) I [posted to Mastodon](https://hachyderm.io/@brainsik/110799062325634606):

[![Black Hole Sun](https://media.hachyderm.io/media_attachments/files/110/766/472/327/295/476/original/2c9c870af2121c9f.png)](https://hachyderm.io/@brainsik/110766510259586543)

[![Cold Wave I](https://media.hachyderm.io/media_attachments/files/110/766/473/972/941/792/original/8790e99669f3d937.png)](https://hachyderm.io/@brainsik/110766510259586543)

[![Cold Wave II](https://media.hachyderm.io/media_attachments/files/110/799/023/645/893/103/original/256903eb291ccd11.png)](https://hachyderm.io/@brainsik/110799062325634606)

## Plane Mapping

* ComplexPoint — A point in the complex plane.
* ImagePoint — A point in the raster image.

The image plane has its origin (0, 0) in the top left corner. The complex plane has it's origin (0+0i) in the center. Increasing `x` and `r` (real) both go to the right. However, increasing `y` and `i` (imaginary) move opposite directions: `x` goes down and `i` goes up.

A "view" of a plane is a rectangle defined by its min and max points. Since `y` and `i` have their min/max points flipped, the view of the image plane is defined by (top-left, bottom-right) while the view of the complex plane is defined by (bottom-left, top-right).

| x | y | <-> | r | i |
|:--:|:--:|:--:|:--:|:--:|
| 0 | 0 |  | -r_max | +i_max |
| 0 | y_max | | -r_max | -i_max |
| x_max | 0 |  | +r_max | +i_max |
| x_max | y_max |  | +r_max | -i_max |
| x_max/2 | y_max/2 |  | 0 | 0 |

## Notes

* [Computer Color is Broken](https://www.youtube.com/watch?v=LKnqECcg6Gw)