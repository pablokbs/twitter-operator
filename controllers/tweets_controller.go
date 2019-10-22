package controllers

import (
	"context"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	webappv1 "twitter/api/v1"
)

// TweetsReconciler reconciles a Tweets object
type TweetsReconciler struct {
	client.Client
	Log       logr.Logger
	TweetDict map[string]int64
	IsInit    int64
}

// +kubebuilder:rbac:groups=webapp.my.domain,resources=tweets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=webapp.my.domain,resources=tweets/status,verbs=get;update;patch

func (r *TweetsReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	//	_ = context.Background()
	//	_ = r.Log.WithValues("tweets", req.NamespacedName)

	// your logic here

	/////////////////////////////////////////twitter credentials

	var (
		consumer_key    = os.Getenv("TWITTER_API_CONSUMER_KEY")
		consumer_secret = os.Getenv("TWITTER_API_CONSUMER_SECRET")
		access_token    = os.Getenv("TWITTER_API_ACCESS_TOKEN")
		access_secret   = os.Getenv("TWITTER_API_ACCESS_TOKEN_SECRET")
	)

	config := oauth1.NewConfig(consumer_key, consumer_secret)
	token := oauth1.NewToken(access_token, access_secret)
	httpClient := config.Client(oauth1.NoContext, token)
	tw_client := twitter.NewClient(httpClient)

	/////////////////////////////////////////////////////////////Init for sync with Twitter
	//	ctx := context.Background()

	if r.IsInit == 0 {
		fmt.Printf("Init synchronization with Twitter\n")
		/////////////////////////////////////////////////////////////Init for sync with Twitter
		TweetNamespace := "default"
		//	ctx := context.Background()

		//	if r.IsInit == 0 {
		r.IsInit = 1
		old_k8sTweetList := &webappv1.TweetsList{}
		//		if err := r.List(context.TODO(), old_k8sTweetList, client.InNamespace(TweetNamespace)); err != nil {
		if err := r.List(context.Background(), old_k8sTweetList); err != nil {
			fmt.Printf("unable to list init tweets---%+v\n", err)
			return ctrl.Result{}, nil
		}

		fmt.Printf("Delete old tweets on k8s cluster--%+v \n", old_k8sTweetList)
		//------1  Delete all tweets for init on the k8s
		for old_k8sCnt, old_k8sTweetItem := range old_k8sTweetList.Items {
			if err := r.Delete(context.TODO(), &old_k8sTweetItem, client.PropagationPolicy(metav1.DeletePropagationBackground)); err != nil {
				fmt.Printf("Failed to delete the old %+vth tweet\n", old_k8sCnt)
			}
			fmt.Printf("track2-old-tweet-Item %+v\n", old_k8sTweetItem.Spec.Message)
		}

		//------2 Get all tweets from the Twitter
		new_TwitterTweetList, _, _ := tw_client.Timelines.UserTimeline(&twitter.UserTimelineParams{
			ScreenName: "pablokbs",
			Count:      20,
		})

		fmt.Printf("track3-tweet-list %+v\n", new_TwitterTweetList)
		//------3 Create tweets on k8s

		autoCreated_k8sTweet := &webappv1.Tweets{}
		fmt.Printf("track4-old-tweet-Item %+v\n", autoCreated_k8sTweet)

		autoCreated_k8sTweetName := "imported-tweet-1"
		for new_k8sCnt, new_TwittertweetItem := range new_TwitterTweetList {
			fmt.Printf("%+v\n", new_TwittertweetItem.ID)
			autoCreated_k8sTweet = &webappv1.Tweets{}
			autoCreated_k8sTweet.ObjectMeta.Name = autoCreated_k8sTweetName + new_TwittertweetItem.IDStr
			autoCreated_k8sTweet.Spec.Message = new_TwittertweetItem.Text
			fmt.Printf("text---%+v\n", new_TwittertweetItem.Text)
			autoCreated_k8sTweet.ObjectMeta.Namespace = TweetNamespace
			autoCreated_k8sTweet.IsAutoCreated = 1
			if err := r.Create(context.TODO(), autoCreated_k8sTweet); err != nil {
				fmt.Printf("Failed to create %+vth k8s auto tweet\n", new_k8sCnt)
				return ctrl.Result{}, nil
			}
			r.TweetDict[autoCreated_k8sTweet.ObjectMeta.Name] = new_TwittertweetItem.ID
		}
		r.IsInit = 1

		return ctrl.Result{}, nil
	}
	/////////////////////////////////////////kubebuilder logic
	fmt.Printf("Reconcile started!\n")
	fmt.Printf("%+v\n", r.TweetDict)

	instance_youma := &webappv1.Tweets{}
	err := r.Get(context.TODO(), req.NamespacedName, instance_youma)

	/////////////////// Delete tweet
	if err != nil {
		/////// Delete tweet on Twitter
		tweet_id := r.TweetDict[req.Name]
		if tweet_id == 0 {
			fmt.Printf("Tweet Id is zero\n")
			return ctrl.Result{}, nil
		}
		fmt.Printf("--Delete tweet %+v\n", tweet_id)

		_, resp_del, err_del := tw_client.Statuses.Destroy(tweet_id, nil)
		fmt.Printf("Delete tweet status code is %+v\n", resp_del.StatusCode)
		if err_del != nil {
			fmt.Printf("Delete tweet API request failed\n")
			return ctrl.Result{}, nil
		}
		////// Delete tweet mapping info
		delete(r.TweetDict, req.Name)
		//fmt.Printf("%+v\n", r.TweetDict)
		return ctrl.Result{}, nil
	}

	//////////////////// Created or modified already

	fmt.Printf("--Create tweet\n")
	//////////// Create tweet on Twitter
	if instance_youma.IsAutoCreated == 1 {
		return ctrl.Result{}, nil
	}
	tweets, resp, _ := tw_client.Statuses.Update(instance_youma.Spec.Message, nil)
	fmt.Printf("%+v\n", resp.StatusCode)
	//fmt.Printf("%+v\n", instance_youma.Status.TweetId)

	if resp.StatusCode != 200 && instance_youma.IsAutoCreated == 0 {
		fmt.Printf("Twitter instance already exists\n")
		if err := r.Delete(context.TODO(), instance_youma, client.PropagationPolicy(metav1.DeletePropagationBackground)); err != nil {
			fmt.Printf("Error while deleting the duplicated instances\n")
		}

		fmt.Printf("Deleted correctly because it exists already\n")
		return ctrl.Result{}, nil
	}
	///////// Create tweet mapping
	if resp.StatusCode == 200 {
		fmt.Printf("--Tweet created\n")
		r.TweetDict[req.Name] = tweets.ID
		//instance_youma.Status.TweetId = tweets.ID
		//err = r.Update(context.TODO(), instance_youma)
		//if err != nil {
		//fmt.Printf("Status Update Failed")
		//return ctrl.Result{}, nil
		//}
	}
	return ctrl.Result{}, nil
}

func (r *TweetsReconciler) SetupWithManager(mgr ctrl.Manager) error {

	fmt.Printf("Watch started!\n")
	r.TweetDict = make(map[string]int64)
	r.IsInit = 0
	retv := ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.Tweets{}).
		Complete(r)
		//	r.Reconcile(ctrl.Request{})
	return retv
}
