package com.ryan3r.caustic.controller;

import com.ryan3r.caustic.Team;
import com.ryan3r.caustic.model.Submission;
import com.ryan3r.caustic.model.accounts;
import com.ryan3r.caustic.repository.LanguageRepository;
import com.ryan3r.caustic.repository.ProblemRepository;
import com.ryan3r.caustic.repository.SubmissionRepository;
import com.ryan3r.caustic.repository.accountsRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.CookieValue;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestParam;

import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.NoSuchElementException;

@Controller
public class PageController {

    @Value("${spring.application.name}")
    String appName;

    private final accountsRepository accRep;

    @Autowired
    private SubmissionRepository submissions;
    @Autowired
    private LanguageRepository languageRepository;
    @Autowired
    private ProblemRepository problemRepository;

    @Autowired
    public PageController(accountsRepository accRep){
        this.accRep = accRep;
    }

    @GetMapping("/")
    public String index(Model model){
        model.addAttribute("appName", appName);
        model.addAttribute("problems", problemRepository.findAll());
        return "index";
    }

    @GetMapping("/profile")
    public String profile(@CookieValue(value="username", required=false) String username, Model model) {

        if (username == null){
            model.addAttribute("error", "1");
            return "profile";
        } else {
            accounts account = accRep.findUser(username);
            model.addAttribute("account", account);
            model.addAttribute("error", "0");
            return "profile";
        }
    }

    @GetMapping("/results/{id}")
    public String results(@PathVariable("id") String id, Model model) {
        Submission submit = null;
        try {
            submit = submissions.findById(Long.parseLong(id)).get();
        } catch(NoSuchElementException ex) {
            model.addAttribute("error", "Submission not found");
        }
        model.addAttribute("result", submit);
        return "results";
    }

    @GetMapping("/scoreboard")
    public String scoreboard(Model model){
        HashMap<String, HashMap<Long, Submission>> subs = new HashMap<>();
        for(Submission sub : submissions.findAll()) {
            if(!subs.containsKey(sub.getSubmitter())) {
                subs.put(sub.getSubmitter(), new HashMap<>());
            }

            HashMap<Long, Submission> userMap = subs.get(sub.getSubmitter());
            if(!userMap.containsKey(sub.getSubmissionId()) || sub.getSolutionTime() < userMap.get(sub.getSubmissionId()).getSolutionTime()) {
                userMap.put(sub.getSubmissionId(), sub);
            }
        }

        ArrayList<Team> teams = new ArrayList<Team>();
        for(String user : subs.keySet()) {
            HashMap<Long, Submission> userMap = subs.get(user);
            Team team = new Team(user, 0);
            teams.add(team);

            for(Long id : userMap.keySet()) {
                team.setScore(team.getScore() + userMap.get(id).getScore());
            }
        }

        Collections.sort(teams);

        int rank = 0;
        for(Team team : teams) {
            team.setRank(++rank);
        }

        model.addAttribute("scores", teams);
        return "scoreboard";
    }

    @GetMapping("/accountSetup")
    public String accountSetup(Model model){
//        AccountDTO userDto = new AccountDTO();
//        model.addAttribute("user", userDto);
        return "accountSetup";
    }

//    @PostMapping("/accountSetup")
//    public ModelAndView registerNewAccount( @ModelAttribute("user") @Valid AccountDTO accountDto,
//                                            BindingResult result,
//                                            WebRequest request,
//                                            Errors errors) {
//
//        accounts registered = new accounts();
//        if (!result.hasErrors()) {
//            registered = createUserAccount(accountDto, result);
//        }
//        if (registered == null) {
//            result.rejectValue("email", "message.regError");
//        }
//        if (result.hasErrors()) {
//            return new ModelAndView("accountSetup", "user", accountDto);
//        }
//        else {
//            return new ModelAndView("formUpload");
//        }
//    }
//    private accounts createUserAccount(AccountDTO accountDto, BindingResult result) {
//        accounts registered = null;
//        try {
//            registered = AccountService.registerNewUserAccount(accountDto);
//        } catch (AccountAlreadyExists e) {
//            return null;
//        }
//        return registered;
//    }

    @GetMapping("/login-newaccount")
    public String login(Model model){
        return "login";
    }

    @GetMapping("/formUpload/{id}")
    public String formUpload(@PathVariable("id") String id, @RequestParam(value = "invalid", required = false) boolean invalid, Model model){
        model.addAttribute("invalid", invalid);
        model.addAttribute("types", languageRepository.findAll());
        model.addAttribute("probId", id);
        model.addAttribute("problemUrl", "/pdf/" + id + ".pdf");
        return "formUpload";
    }
}