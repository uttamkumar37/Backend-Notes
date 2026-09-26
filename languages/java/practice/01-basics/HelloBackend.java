public class HelloBackend {
    public static void main(String[] args) {
        String language = "Java";
        String goal = "backend engineering";

        printWelcome(language, goal);
        printNextSteps();
    }

    private static void printWelcome(String language, String goal) {
        System.out.println("Learning " + language + " for " + goal + ".");
    }

    private static void printNextSteps() {
        System.out.println("Next: practice variables, methods, classes, and collections.");
    }
}

