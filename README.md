## ToDo

* BUG: WTF is going on with ap.calc_area?? Something is not right. Changing the size of the plane seems change where you're oriented in it. This is perhaps a problem with MakeProblemSet because adding a separate MakeProblemSet for Juliabrots (where we iterate over the image plane instead of the complex plane to make the problem set) makes the issue go away. Something is up.
* If we fix this bug ☝️ then we _may_ be able to drop the two MakeProblemSet functions and have just one. I say may because the z <-> xy mapping is not 1:1 and separating problem set creators by whether they are based on the complex plane or the image plane is maybe just correct.
* ?? Instead of PlanePoint and ImagePoint, move to CmplxPoint and ImagePoint and rename other things accordingly?
* Add godoc comments to everything
* Use Log() instead of Printf()?
* ColorFunc: Add a params type so we can parameterize them.

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