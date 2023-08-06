## ToDo

* Use Log() instead of Printf()?
* ColorFunc: Add a params type so we can parameterize them.
* Nuke CalcPoint type? We aren't using the xy component at all … do we need this or can we use a basic type instead of the struct?

## Plane Mapping

* PlanePoint — A point in the complex plane.
* ImagePoint — A point in the image.

The image plane has its origin (0, 0) in the top left corner. The complex plane has it's origin (0+0i) in the center. Increase `x` and `r` both go to the right. However, increasing `y` and `i` move opposite directions: `x` goes down and `i` goes up.

A "view" of a plane is a rectangle defined by its min and max points. Since `y` and `i` have their min/max points flipped, the view of the image plane is defined by (top-left, bottom-right) while the view of the complex plane is defined by (bottom-left, top-right).

| x | y | <-> | r | i |
|:--:|:--:|:--:|:--:|:--:|
| 0 | 0 |  | -r_max | +i_max |
| 0 | y_max | | -r_max | -i_max |
| x_max | 0 |  | +r_max | +i_max |
| x_max | y_max |  | +r_max | -i_max |
| x_max/2 | y_max/2 |  | 0 | 0 |

## Notes

[Computer Color is Broken](https://www.youtube.com/watch?v=LKnqECcg6Gw)