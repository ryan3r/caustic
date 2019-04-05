package com.ryan3r.caustic.controller;

import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;

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

        return "profile";
    }

    @GetMapping("/results")
    public String results(Model model, @RequestParam(value="result", required=false, defaultValue="fail") String result) {
        model.addAttribute("result", result);
        return "results";
    }

    @GetMapping("/scoreboard")
    public String scoreboard(Model model){
        return "scoreboard";
    }
}