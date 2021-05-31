package utils;
import java.awt.Point;

public class Utils {
	public static boolean arePointsEqual(Point p1, Point p2) {	
		if (p1.x != p2.x || p1.y != p2.y)
			return false;
		
		return true;	
	}
}
