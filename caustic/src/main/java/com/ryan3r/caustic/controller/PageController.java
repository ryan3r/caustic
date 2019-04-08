package com.ryan3r.caustic.controller;

import com.ryan3r.caustic.model.Scoreboard;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;

import java.util.ArrayList;
import java.util.List;

@Controller
public class PageController {

//    @GetMapping(path="/")
//    public String index(Model model){
//        return "index";
//    }

    @GetMapping("/profile")
    public String profile(Model model,
                          @RequestParam(value="firstName", required=false, defaultValue="Zach") String firstName,
                          @RequestParam(value="lastName", required=false, defaultValue="Gorman") String lastName,
                          @RequestParam(value="email", required=false, defaultValue="zgorman2@iastate.edu") String email,
                          @RequestParam(value="rank", required=false, defaultValue="") int rank) {

        model.addAttribute("firstName", firstName);
        model.addAttribute("lastName", lastName);
        model.addAttribute("email", email);
        model.addAttribute("rank", rank);

        return "profile.html";
    }

    @GetMapping("/results")
    public String results(Model model, @RequestParam(value="result", required=false, defaultValue="fail") String result) {
        model.addAttribute("result", result);
        return "results.html";
    }

    List<Scoreboard> scoreList = new ArrayList<>();

    @GetMapping("/scoreboard")
    public String scoreboard(Model model){
        model.addAttribute("scores", scoreList);
        return "scoreboard.html";
    }
}