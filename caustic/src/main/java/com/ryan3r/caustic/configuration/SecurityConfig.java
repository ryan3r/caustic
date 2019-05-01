//package com.ryan3r.caustic.configuration;
//
//import org.springframework.beans.factory.annotation.Autowired;
//import org.springframework.context.annotation.Bean;
//import org.springframework.context.annotation.Configuration;
//import org.springframework.security.authentication.dao.DaoAuthenticationProvider;
//import org.springframework.security.config.annotation.authentication.builders.AuthenticationManagerBuilder;
//import org.springframework.security.config.annotation.web.builders.HttpSecurity;
//import org.springframework.security.config.annotation.web.builders.WebSecurity;
//import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
//import org.springframework.security.config.annotation.web.configuration.WebSecurityConfigurerAdapter;
//import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
//import org.springframework.security.crypto.password.PasswordEncoder;
//import org.springframework.web.context.WebApplicationContext;
//import org.springframework.security.data.repository.query.SecurityEvaluationContextExtension;
//
//
//import javax.annotation.PostConstruct;
//import javax.sql.DataSource;
//
//@Configuration
//@EnableWebSecurity
//public class SecurityConfig extends WebSecurityConfigurerAdapter {
//
//    @Autowired
//    private DataSource dataSource;
//
//    @Autowired
//    private WebApplicationContext applicationContext;
//
//    private CAUSTICUserDetailsService usrDetailServ;
//
//    public SecurityConfig() {
//        super();
//    }
//
//    @Override
//    protected void configure(AuthenticationManagerBuilder auth) throws Exception {
//        auth
//            .userDetailsService(usrDetailServ)
//                .passwordEncoder(encoder())
//                .and()
//            .authenticationProvider(authenticationProvider())
//            .jdbcAuthentication()
//            .usersByUsernameQuery("select * from accounts account where account.username=@?")
//            .dataSource(dataSource);
//    }
//
//    @Override
//    protected void configure(final HttpSecurity http) throws Exception {
//        http
//            .authorizeRequests()
//                .antMatchers("/", "/accountSetup", "/error", "/scoreboard").permitAll()
//                .anyRequest().authenticated()
//                .and()
//            .formLogin()
//                .loginPage("/login-newaccount").permitAll()
//                .failureUrl("/login-newaccount?error")
//                .and()
//            .csrf().disable();
//    }
//
//    @Override
//    public void configure(WebSecurity web) {
//        web
//            .ignoring()
//            .antMatchers("/resources/**", "/static/**");
//    }
//
//    @Bean
//    public DaoAuthenticationProvider authenticationProvider() {
//
//        final DaoAuthenticationProvider authProvider = new DaoAuthenticationProvider();
//
//        authProvider.setUserDetailsService(usrDetailServ);
//        authProvider.setPasswordEncoder(encoder());
//
//        return authProvider;
//    }
//
//    @PostConstruct
//    public void completeSetup() {
//
//        usrDetailServ = applicationContext.getBean(CAUSTICUserDetailsService.class);
//    }
//
//    @Bean
//    public SecurityEvaluationContextExtension securityEvaluationContextExtension() {
//        return new SecurityEvaluationContextExtension();
//    }
//
//    @Bean
//    public PasswordEncoder encoder() {
//        return new BCryptPasswordEncoder();
//    }
//}
