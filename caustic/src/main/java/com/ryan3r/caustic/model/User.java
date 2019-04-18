package com.ryan3r.caustic.model;

import javax.persistence.Entity;
import javax.persistence.Id;

@Entity
public class User {

    @Id
    private String name;
    private int rank;
    private String email;

    public User(){}

    public User(String name, String email){
        this.name = name;
        this.email = email;
    }
}
