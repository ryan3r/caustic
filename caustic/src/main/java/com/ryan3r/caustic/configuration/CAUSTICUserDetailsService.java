//package com.ryan3r.caustic.configuration;
//
//import com.ryan3r.caustic.model.accounts;
//import com.ryan3r.caustic.repository.accountsRepository;
//import org.springframework.beans.factory.annotation.Autowired;
//import org.springframework.security.core.userdetails.UserDetails;
//import org.springframework.security.core.userdetails.UserDetailsService;
//import org.springframework.security.core.userdetails.UsernameNotFoundException;
//import org.springframework.stereotype.Service;
//import org.springframework.web.context.WebApplicationContext;
//
//import javax.annotation.PostConstruct;
//
//@Service
//public class CAUSTICUserDetailsService implements UserDetailsService {
//
//    @Autowired
//    private WebApplicationContext webAppCtxt;
//    private accountsRepository a;
//
//    public CAUSTICUserDetailsService(){
//        super();
//    }
//
//    @Override
//    public UserDetails loadUserByUsername(String username) throws UsernameNotFoundException {
//
//        accounts account = a.findUser(username);
//        if (account == null){
//            throw new UsernameNotFoundException(username);
//        }
//
//        return new AccountPrinciple(account);
//    }
//
//    @PostConstruct
//    public void completeSetup() {
//        a = webAppCtxt.getBean(accountsRepository.class);
//    }
//}
