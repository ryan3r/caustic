package main.java.com.ryan3r.caustic.model;

public class User {

    private String firstName;
    private String lastName;
    private int rank;
    private String email;

    public User(){}

    public User(String firstName, String lastName, String email){
        this.firstName = firstName;
        this.lastName = lastName;
        this.email = email;
    }
}
