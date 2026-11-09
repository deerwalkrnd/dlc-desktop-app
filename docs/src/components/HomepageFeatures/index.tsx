import type {ReactNode} from 'react';
import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

type FeatureItem = {
  title: string;
  description: ReactNode;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Works offline',
    description: (
      <>
        One small program and a folder of videos. No internet connection, no
        installation, no database.
      </>
    ),
  },
  {
    title: 'Organised automatically',
    description: (
      <>
        Videos are sorted into classes, subjects, chapters and lessons straight
        from their file names. Add files, press <strong>Rescan</strong>, done.
      </>
    ),
  },
  {
    title: 'Built for classrooms',
    description: (
      <>
        A clear player with keyboard shortcuts, playback speed, full screen and
        next-lesson navigation, in any modern browser on Windows 10 or later.
      </>
    ),
  },
];

function Feature({title, description}: FeatureItem) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center padding-horiz--md">
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): ReactNode {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props) => (
            <Feature key={props.title} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
