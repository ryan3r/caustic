package main.java.com.zach.caustic.model;

import java.io.File;

public class Results {

    private File code;
    private String name;
    private int score;
    private String passFail;

    public Results(){}

    public Results(File code, String name){
        this.code = code;
        this.name = name;
    }
}
