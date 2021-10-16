use std::collections::HashMap;

#[derive(Debug)]
struct Rules {
    // e.g. scca, nasa, champcar, etc
    organization: String,
    classes: HashMap<String, Class>
}

#[derive(Debug)]
struct Class {
    // e.g. Street, StreetTouring, StreetModified, etc
    name: String,
    subclasses: HashMap<String, bool>,
    bump_questions: HashMap<String, BumpQuestion>
}

// BumpQuestions are just yes/no questions that will bump someone out of a class if they
// answered "no"
#[derive(Debug)]
struct BumpQuestion {
    question_prompt: String,
    question_body: String
}

fn main() {
    let street_fender_question: BumpQuestion = BumpQuestion {
        question_prompt: String::from("Are your fenders unmodified?"),
        question_body: String::from("")
    };

    let street_tire_question: BumpQuestion = BumpQuestion {
        question_prompt: String::from("Are your tires 200 treadwear and DOT legal?"),
        question_body: String::from("")
    };

    let mut street_subclasses = HashMap::new();

    street_subclasses.insert(String::from("A Street (AS)"), true);
    street_subclasses.insert(String::from("B Street (BS)"), true);

    let mut street_bump_questions = HashMap::new();

    street_bump_questions.insert(String::from("Fenders"), street_fender_question);
    street_bump_questions.insert(String::from("Tires"), street_tire_question);

    let street_class = Class {
        name: String::from("Street"),
        subclasses: street_subclasses,
        bump_questions: street_bump_questions
    };

    let mut classes = HashMap::new();

    classes.insert(String::from("street"), street_class);

    let scca_rules = Rules {
        organization: String::from("SCCA"),
        classes: classes
    };

    println!("{:?}", scca_rules);

    populate_rules()
}

fn populate_rules() {
    println!("ohai");
}
