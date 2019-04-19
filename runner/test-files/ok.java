import java.util.*;
import java.io.*;

public class ok {
    public static void main(String[] args) {
        Scanner s = new Scanner(System.in);
        while(s.hasNextInt()) System.out.print(s.nextInt() + " ");
        System.out.println();
        s.close();
    }
}