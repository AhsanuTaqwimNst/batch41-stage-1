const dataProject = [];

function addProject(event) {
  event.preventDefault();

  let projectTitle = document.getElementById("project-name").value;
  let startDate = document.getElementById("start-date").value;
  let endDate = document.getElementById("end-date").value;
  let projectDescription = document.getElementById("project-description").value;
  let useJava = document.getElementById("use-java").checked;
  let useNodeJS = document.getElementById("use-node-js").checked;
  let useReact = document.getElementById("use-react").checked;
  let useFacebook = document.getElementById("use-facebook").checked;

  let uploadImage = document.getElementById("upload-image").files[0];

  image = URL.createObjectURL(uploadImage);

  let project = {
    projectTitle,
    startDate,
    endDate,
    projectDescription,
    useJava,
    useNodeJS,
    useReact,
    useFacebook,
    image,
  };
  dataProject.push(project);
  console.log(dataProject);

  renderProject();
}

function renderProject() {
  let projectList = document.getElementById("project-content");

  projectList.innerHTML = "";

  for (let index = 0; index < dataProject.length; index++) {
    projectList.innerHTML += `

      <div class="project-item">
        <a href="blog-project.html"
          ><img src="${dataProject[index].image}" alt=""
        /></a>
        <div class="project-detail">
          <h3>${dataProject[index].projectTitle}</h3>
          <p>Duration: ${getDurationTime(
            dataProject[index].startDate,
            dataProject[index].endDate
          )}</p>
        </div>
        <div class="project-description">
          <p>${dataProject[index].projectDescription}</p>
        </div>
        <div class="project-tech">
            ${
              dataProject[index].useJava
                ? `<i class="fa-brands fa-java"></i>`
                : ""
            }
            ${
              dataProject[index].useNodeJS
                ? `<i class="fa-brands fa-node-js"></i>`
                : ""
            }
            ${
              dataProject[index].useReact
                ? `<i class="fa-brands fa-react"></i>`
                : ""
            }
            ${
              dataProject[index].useFacebook
                ? `<i class="fa-brands fa-facebook"></i>`
                : ""
            }
        </div>
        <div class="action-btn">
          <button class="btn">edit</button>
          <button class="btn">delete</button>
        </div>
      </div> 
    `;
  }
}

function getDurationTime(start, end) {
  let dateStart = new Date(start);
  let dateEnd = new Date(end);

  let diffTime = dateEnd - dateStart;

  let distanceSecond = Math.floor(diffTime / 1000);
  let distanceMinute = Math.floor(diffTime / (1000 * 60));
  let distanceHour = Math.floor(diffTime / (1000 * 60 * 60));
  let distanceDay = Math.floor(diffTime / (1000 * 60 * 60 * 24));
  let distanceWeek = Math.floor(diffTime / (1000 * 60 * 60 * 24 * 7));
  let distanceMonth = Math.floor(diffTime / (1000 * 60 * 60 * 24 * 30));
  let distanceYear = Math.floor(diffTime / (1000 * 60 * 60 * 24 * 30 * 12));

  if (distanceYear > 0) {
    return `${distanceYear} year(s) ago`;
  } else if (distanceMonth > 0) {
    return `${distanceMonth} month(s) ago`;
  } else if (distanceWeek > 0) {
    return `${distanceWeek} week(s) ago`;
  } else if (distanceDay > 0) {
    return `${distanceDay} day(s) ago`;
  } else if (distanceHour > 0) {
    return `${distanceHour} hour(s) ago`;
  } else if (distanceMinute > 0) {
    return `${distanceMinute} minute(s) ago`;
  } else if (distanceSecond > 0) {
    return `${distanceSecond} second(s) ago`;
  }
}
