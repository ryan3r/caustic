package com.ryan3r.caustic.model;

public class Scoreboard {

    private int rank;
    private String teamName;
    private boolean solved;
    private int timeToSolve;
    private double score;

    public Scoreboard(){}

    public Scoreboard(int rank, String teamName, boolean solved, int timeToSolve, double score){
        this.rank = rank;
        this.teamName = teamName;
        this.solved = solved;
        this.timeToSolve = timeToSolve;
        this.score = score;
    }
}
