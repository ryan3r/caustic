package com.ryan3r.caustic.model;

import javax.persistence.Entity;

@Entity
public class User {

    private String name;
    private int rank;
    private String email;

    public User(){}

    public User(String name, String email){
        this.name = name;
        this.email = email;
    }
}
