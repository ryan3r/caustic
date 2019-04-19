package com.ryan3r.caustic.repository;

import com.ryan3r.caustic.model.accounts;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.CrudRepository;
import org.springframework.data.repository.query.Param;

public interface accountsRepository extends CrudRepository<accounts, String>{
    @Query(value = "SELECT * FROM accounts account WHERE account.username = :username", nativeQuery = true)
    accounts findUser(@Param("username") String username);
}
