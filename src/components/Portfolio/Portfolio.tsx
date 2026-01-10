import { Suspense } from "react";
import { useFontsLoaded } from "../hooks/useFontsLoaded";

import PortfolioItems, {
  type PortfolioItem,
} from "@components/Headless/Portfolio/PortfolioItems";

interface CaseStudy {
  id: string;
  title: string;
  year: string;
  subtitle: string;
  description: string;
  keyAchievements?: string[];
  portfolioItems: PortfolioItem[];
}

const caseStudies: CaseStudy[] = [
  {
    id: "hailstorm-talent",
    title: "HailStorm Talent",
    year: "2025-Present",
    subtitle: "Talent Onboarding & Contract Management System",
    description:
      "Streamlined the talent onboarding workflow for union and non-union projects by implementing an end-to-end document management system.",
    keyAchievements: [
      "PDF Contract Generation: Automated contract creation workflow with dynamic templates for different talent types and project requirements",
      "E-Signature Integration: Implemented email delivery system for contracts to talent and producers, integrated digital signature capture and verification",
      "Document Archival: Built secure document storage and retrieval system for signed contracts and onboarding materials",
      "Tax & Employment Forms: Created workflow for I-9, W-4, and union-specific tax forms with validation and compliance checking",
      "Onboarding Status Tracking: Developed real-time dashboard to monitor talent onboarding progress from initial contact through full clearance",
    ],
    portfolioItems: [
      {
        src: "/screenshots/hailstorm/iPad Landscape.avif",
        alt: "HailStorm Talent Dashboard",
      },
      {
        src: "/screenshots/hailstorm/iPad Landscape 2.avif",
        alt: "HailStorm Contract Management",
      },
      {
        src: "/screenshots/hailstorm/iPad Landscape 3.avif",
        alt: "HailStorm Document Workflow",
      },
      {
        src: "/screenshots/hailstorm/iPad Landscape 4.avif",
        alt: "HailStorm Onboarding Status",
      },
    ],
  },
  {
    id: "attesto",
    title: "Attesto",
    year: "2024",
    subtitle: "Front End Demo for clients",
    description:
      "Converted Figma designs into React components with interactivity to mimic a real workflow.",
    portfolioItems: [
      {
        src: "/screenshots/attesto/attesto-01-checkin.png",
        alt: "Attesto Check-In",
      },
      {
        src: "/screenshots/attesto/attesto-02-welcome.png",
        alt: "Attesto Welcome",
      },
      {
        src: "/screenshots/attesto/attesto-03-role.png",
        alt: "Attesto Select Role",
      },
      {
        src: "/screenshots/attesto/attesto-04-questions.png",
        alt: "Attesto Questions",
      },
      {
        src: "/screenshots/attesto/attesto-05-add-question.png",
        alt: "Attesto Add Question",
      },
      {
        src: "/screenshots/attesto/attesto-06-select-candidates.png",
        alt: "Attesto Select Candidates",
      },
      {
        src: "/screenshots/attesto/attesto-07-upload-resume.png",
        alt: "Attesto Upload Resume",
      },
      {
        src: "/screenshots/attesto/attesto-08-check-assessments.png",
        alt: "Attesto Check Assessments",
      },
      {
        src: "/screenshots/attesto/attesto-09-assignment-popup.png",
        alt: "Attesto Assignment Popup",
      },
      {
        src: "/screenshots/attesto/attesto-10-questions.png",
        alt: "Attesto Assignment Questions",
      },
      {
        src: "/screenshots/attesto/attesto-11-questions-qa.png",
        alt: "Attesto Questions QA",
      },
      {
        src: "/screenshots/attesto/attesto-12-end-demo.png",
        alt: "Attesto End Demo",
      },
    ],
  },
];

export default function Portfolio() {
  const fontsLoaded = useFontsLoaded();

  if (!fontsLoaded) {
    return <div>Loading...</div>;
  }

  return (
    <div className="flex flex-col space-y-12">
      <h1 className="text-3xl text:xs-3xl md:text-6xl font-semibold text-white">
        Projects
      </h1>

      {caseStudies.map((caseStudy, index) => (
        <section key={caseStudy.id} className="space-y-4">
          <div className="border-t border-gray-700 pt-8">
            <h2 className="text-2xl md:text-4xl font-semibold">
              {caseStudy.title} - {caseStudy.year}
            </h2>
            <h3 className="text-xl md:text-2xl text-gray-300 mt-2">
              {caseStudy.subtitle}
            </h3>
            <p className="text-gray-400 mt-4">{caseStudy.description}</p>

            {caseStudy.keyAchievements &&
              caseStudy.keyAchievements.length > 0 && (
                <div className="mt-6">
                  <h4 className="text-lg md:text-xl font-semibold mb-3">
                    Key Contributions:
                  </h4>
                  <ul className="space-y-2">
                    {caseStudy.keyAchievements.map((achievement, i) => (
                      <li
                        key={i}
                        className="text-gray-400 pl-4 border-l-2 border-blue-500"
                      >
                        {achievement}
                      </li>
                    ))}
                  </ul>
                </div>
              )}
          </div>

          <PortfolioItems items={caseStudy.portfolioItems} />
        </section>
      ))}
    </div>
  );
}
