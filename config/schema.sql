create table if not exists apps
(
    bundle           String,
    developerId      String,
    developer        String,
    title            String,
    categories       String,
    price            String,
    picture          String,
    screenshots      Array(String),
    rating           String,
    reviewCount      String,
    ratingHistogram  Array(String),
    description      String,
    shortDescription String,
    recentChanges    String,
    releaseDate      String,
    lastUpdateDate   String,
    appSize          String,
    installs         String,
    version          String,
    androidVersion   String,
    contentRating    String,
    developerContacts Nested(email String, contacts String),
    privacyPolicy    String,
    datetime         DateTime DEFAULT now()
)   ENGINE = ReplacingMergeTree(datetime)
    ORDER BY (bundle, developerId, categories)
    PARTITION BY (categories);