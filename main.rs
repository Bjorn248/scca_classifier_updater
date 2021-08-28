#[derive(Debug)]
struct Rules {
    // e.g. scca, nasa, champcar, etc
    organization: String,
    classes: Vec<Class>
}

#[derive(Debug)]
struct Class {
    // e.g. Street, StreetTouring, StreetModified, etc
    name: String,
    subclasses: Vec<SubClass>,
    bump_questions: Vec<BumpQuestion>
}

#[derive(Debug)]
struct SubClass {
    // e.g. BS, STU, SMF, etc
    name: String
}

// BumpQuestions are just yes/no questions that will bump someone out of a class if they
// answered "no"
#[derive(Debug)]
struct BumpQuestion {
    question_prompt: String,
    question_body: String
}

fn main() {
    let street_fender_question = BumpQuestion {
        question_prompt: String::from("Are your fenders unmodified?"),
        question_body: String::from("")
    };

    let street_tire_question = BumpQuestion {
        question_prompt: String::from("Are your tires 200 treadwear and DOT legal?"),
        question_body: String::from("")
    };

    let a_street = SubClass {
        name: String::from("A Street (AS)")
    };

    let b_street = SubClass {
        name: String::from("B Street (BS)")
    };

    let street_subclasses = vec![a_street, b_street];

    let street_bump_questions = vec![street_fender_question, street_tire_question];

    let street = Class {
        name: String::from("Street"),
        subclasses: street_subclasses,
        bump_questions: street_bump_questions
    };

    let classes = vec![street];

    let scca_rules = Rules {
        organization: String::from("SCCA"),
        classes: classes
    };

    println!("hi :3");
    println!("{:?}", scca_rules);
}
