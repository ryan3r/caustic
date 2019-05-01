//package com.ryan3r.caustic.configuration;
//
//import com.ryan3r.caustic.model.accounts;
//import org.springframework.security.core.GrantedAuthority;
//import org.springframework.security.core.authority.SimpleGrantedAuthority;
//import org.springframework.security.core.userdetails.UserDetails;
//
//import java.util.Collection;
//import java.util.Collections;
//
//public class AccountPrinciple implements UserDetails {
//
//    private final accounts acc;
//
//    public AccountPrinciple(accounts account){
//        acc = account;
//    }
//
//    @Override
//    public Collection<? extends GrantedAuthority> getAuthorities() {
//        return Collections.<GrantedAuthority>singletonList(new SimpleGrantedAuthority("comp"));
//    }
//
//    @Override
//    public String getUsername() {
//        return acc.getUsername();
//    }
//
//    @Override
//    public String getPassword() {
//        return acc.getPassword();
//    }
//
//    @Override
//    public boolean isAccountNonExpired() {
//        return true;
//    }
//
//    @Override
//    public boolean isAccountNonLocked() {
//        return true;
//    }
//
//    @Override
//    public boolean isCredentialsNonExpired() {
//        return true;
//    }
//
//    @Override
//    public boolean isEnabled() {
//        return true;
//    }
//}
