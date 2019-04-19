package com.ryan3r.caustic;

import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.CrudRepository;
import org.springframework.data.repository.query.Param;

public interface accountsRepository extends CrudRepository<accounts, String>{
    @Query("SELECT account FROM accounts WHERE account.username = :username")
    accounts findUser(@Param("username") String username);
}
